package server

type Reaction struct {
	ID        int    `json:"id"`
	MessageID int    `json:"message_id"`
	UserID    int    `json:"user_id"`
	Emoji     string `json:"emoji"`
	Username  string `json:"username"`
}

type ReactionGroup struct {
	Emoji  string   `json:"emoji"`
	Count  int      `json:"count"`
	UserID int      `json:"user_id"`
	Users  []string `json:"users"`
}

func GetReactions(messageID int, currentUserID int) ([]ReactionGroup, error) {
	rows, err := dbConn.Query(`
		SELECT emoji, COUNT(*) as cnt,
		       ARRAY_AGG(username) as usernames,
		       BOOL_OR(user_id = $2) as is_mine
		FROM reactions r
		JOIN users u ON u.id = r.user_id
		WHERE r.message_id = $1
		GROUP BY emoji
		ORDER BY cnt DESC
	`, messageID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ReactionGroup

	for rows.Next() {
		var g ReactionGroup
		var usernames []string
		var isMine bool

		if err := rows.Scan(&g.Emoji, &g.Count, &usernames, &isMine); err != nil {
			return nil, err
		}

		g.Users = usernames
		if isMine {
			g.UserID = currentUserID
		}

		groups = append(groups, g)
	}

	return groups, rows.Err()
}

func AddReaction(messageID int, userID int, emoji string) error {
	_, err := dbConn.Exec(`
		INSERT INTO reactions (message_id, user_id, emoji)
		VALUES ($1, $2, $3)
		ON CONFLICT (message_id, user_id, emoji) DO NOTHING
	`, messageID, userID, emoji)
	return err
}

func RemoveReaction(messageID int, userID int, emoji string) error {
	_, err := dbConn.Exec(`
		DELETE FROM reactions
		WHERE message_id = $1 AND user_id = $2 AND emoji = $3
	`, messageID, userID, emoji)
	return err
}
