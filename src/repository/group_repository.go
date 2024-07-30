package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (repo *GroupRepository) CreateGroup(group *models.Group) error {
	query := `INSERT INTO groups (name) VALUES (?)`
	res, err := repo.db.Exec(query, group.Name)
	if err != nil {
		return err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	group.ID = groupID
	for _, member := range group.Members {
		_, err := repo.db.Exec(`INSERT INTO group_members (group_id, member) VALUES (?, ?)`, group.ID, member)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *GroupRepository) GetGroupMembers(groupID int64) ([]string, error) {
	var members []string
	query := `SELECT member FROM group_members WHERE group_id = ?`
	err := repo.db.Select(&members, query, groupID)
	return members, err
}

func (repo *GroupRepository) SaveGroupMessage(message *models.GroupMessage) error {
	query := `INSERT INTO group_messages (group_id, sender, content, timestamp) VALUES (?, ?, ?, ?)`
	_, err := repo.db.Exec(query, message.GroupID, message.Sender, message.Content, message.Timestamp)
	return err
}

func (repo *GroupRepository) GetGroupMessages(groupID int64) ([]models.GroupMessage, error) {
	var messages []models.GroupMessage
	query := `SELECT * FROM group_messages WHERE group_id = ? ORDER BY timestamp`
	err := repo.db.Select(&messages, query, groupID)
	return messages, err
}
