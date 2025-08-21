package services

import (
	"database/sql"
	"fmt"

	"ergracer-api/internal/models"
)

type FriendshipService struct {
	db *sql.DB
}

func NewFriendshipService(db *sql.DB) *FriendshipService {
	return &FriendshipService{db: db}
}

func (s *FriendshipService) CanInviteFriend(userID, friendID int) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM race_participants rp1
		JOIN race_participants rp2 ON rp1.race_id = rp2.race_id
		WHERE rp1.user_id = $1 AND rp2.user_id = $2`
	
	var hasSharedRace bool
	err := s.db.QueryRow(query, userID, friendID).Scan(&hasSharedRace)
	if err != nil {
		return false, err
	}

	return hasSharedRace, nil
}

func (s *FriendshipService) InviteFriend(userID, friendID int) error {
	canInvite, err := s.CanInviteFriend(userID, friendID)
	if err != nil {
		return err
	}

	if !canInvite {
		return fmt.Errorf("you must participate in at least one race together before sending a friend request")
	}

	query := `
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (user_id, friend_id) DO NOTHING`
	
	_, err = s.db.Exec(query, userID, friendID)
	return err
}

func (s *FriendshipService) AcceptFriendship(userID, friendID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE friendships 
		SET status = 'accepted', accepted_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'`
	
	result, err := tx.Exec(query, friendID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no pending friendship request found")
	}

	reverseQuery := `
		INSERT INTO friendships (user_id, friend_id, status, accepted_at)
		VALUES ($1, $2, 'accepted', CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, friend_id) DO UPDATE SET
			status = 'accepted',
			accepted_at = CURRENT_TIMESTAMP`
	
	_, err = tx.Exec(reverseQuery, userID, friendID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *FriendshipService) GetFriends(userID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.email_verified, u.created_at, u.updated_at
		FROM users u
		JOIN friendships f ON u.id = f.friend_id
		WHERE f.user_id = $1 AND f.status = 'accepted'`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var friend models.User
		err := rows.Scan(
			&friend.ID, &friend.Email, &friend.Username, 
			&friend.EmailVerified, &friend.CreatedAt, &friend.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (s *FriendshipService) GetPendingInvitations(userID int) ([]models.Friendship, error) {
	query := `
		SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.accepted_at
		FROM friendships f
		WHERE f.friend_id = $1 AND f.status = 'pending'`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []models.Friendship
	for rows.Next() {
		var invitation models.Friendship
		err := rows.Scan(
			&invitation.ID, &invitation.UserID, &invitation.FriendID,
			&invitation.Status, &invitation.CreatedAt, &invitation.AcceptedAt,
		)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}