package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type EntityType int

const (
	Server EntityType = iota
	Channel
	All
)

type Entity struct {
	ID              string
	Type            EntityType
	ReceiveChannels []int
	SendChannels    []int
	CreatedAt       time.Time
}

func FetchEntity(id string) (*Entity, error) {
	row := DBConnection.QueryRow("SELECT * FROM `relay_entities` WHERE `id` = ?", id)

	var (
		entity          = &Entity{}
		receiveChannels string
		sendChannels    string
	)

	err := row.Scan(
		&entity.ID,
		&entity.Type,
		&receiveChannels,
		&sendChannels,
		&entity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	entity.ReceiveChannels = ParseChannels(receiveChannels)

	entity.SendChannels = ParseChannels(sendChannels)

	return entity, nil
}

func FetchEntities(eType EntityType) ([]*Entity, error) {
	rows, err := DBConnection.Query("SELECT * FROM `relay_entities` WHERE `type` = ?", eType.Polarize())

	if err != nil {
		return nil, err
	}

	var entities []*Entity

	defer rows.Close()

	for rows.Next() {
		var (
			entity          = &Entity{}
			receiveChannels string
			sendChannels    string
		)

		rows.Scan(
			&entity.ID,
			&entity.Type,
			&receiveChannels,
			&sendChannels,
			&entity.CreatedAt,
		)

		entity.ReceiveChannels = ParseChannels(receiveChannels)

		entity.SendChannels = ParseChannels(sendChannels)

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (entity *Entity) UpdateChannels() (sql.Result, error) {
	return DBConnection.Exec(
		"UPDATE `relay_entities` SET `receive_channels` = ?, `send_channels` = ? WHERE `id` = ?",
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
		entity.ID,
	)
}

func (entity *Entity) CreateEntity() (sql.Result, error) {
	return DBConnection.Exec(
		"INSERT INTO `relay_entities` (`id`, `type`, `receive_channels`, `send_channels`) VALUES (?, ?, ?, ?)",
		entity.ID,
		entity.Type,
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
	)
}

func ParseChannels(s string) (c []int) {
	ss := strings.Split(strings.Replace(s, " ", "", -1), ",")

	for _, channel := range ss {
		tc, _ := strconv.Atoi(channel)
		c = append(c, tc)
	}

	return
}

func EncodeChannels(channels []int) string {
	var s []string

	for _, c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	return strings.Join(s, ",")
}

func ChannelString(channels []int) string {
	var s []string

	for _, c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	j := strings.Join(s, ", ")

	if j == "" {
		return "None"
	}

	return j
}

func (eType EntityType) Polarize() EntityType {
	switch eType {
	case Server:
		return Channel
	case Channel:
		return Server
	default:
		return All
	}
}

func (eType EntityType) String() string {
	switch eType {
	case Server:
		return "Server"
	case Channel:
		return "Channel"
	default:
		return "Unknown"
	}
}