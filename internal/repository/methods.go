package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pr_service/internal/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) CreateTeam(ctx context.Context, team model.Team, users []*model.User) (*int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %v", err)
	}

	qTeam, aTeam, err := sq.
		Insert(TeamsTable).
		Columns(TeamsColumns...).
		PlaceholderFormat(sq.Dollar).
		Values(team.Name).
		Suffix("ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id").
		ToSql()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("build create team sql: %v", err)
	}

	var teamID int64
	if err := tx.QueryRowContext(ctx, qTeam, aTeam...).Scan(&teamID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("exec create team: %v", err)
	}

	build := sq.
		Insert(UsersTable).
		Columns(UsersColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, u := range users {
		u.TeamId = teamID
		build = build.Values(u.Values()...)
	}

	build = build.Suffix(`
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            is_active = EXCLUDED.is_active,
            team_id = EXCLUDED.team_id
    `)

	qUsers, aUsers, err := build.ToSql()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("build create user sql: %v", err)
	}

	if _, err := tx.ExecContext(ctx, qUsers, aUsers...); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("exec create user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %v", err)
	}

	return &teamID, nil
}

func (r *UserRepository) CreatePullRequest(ctx context.Context, pr model.PullRequest) error {
	qTeam, aTeam, err := sq.
		Select("team_id").
		From(UsersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": pr.AuthorId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build find team sql: %v", err)
	}

	var teamID int64
	if err := r.db.QueryRowContext(ctx, qTeam, aTeam...).Scan(&teamID); err != nil {
		return fmt.Errorf("exec find team: %v", err)
	}

	qReviewers, aReviewers, err := sq.
		Select("id").
		From(UsersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.NotEq{"id": pr.AuthorId}, 
			sq.Eq{"team_id": teamID}, 
			sq.Eq{"is_active": true}}).
		Limit(2).
		ToSql()
	if err != nil {
		return fmt.Errorf("build reviewers sql: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, qReviewers, aReviewers...)
	if err != nil {
		return fmt.Errorf("exec reviewers query: %v", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan reviewer: %v", err)
		}
		reviewers = append(reviewers, id)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows err: %v", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %v", err)
	}

	qPr, aPr, err := sq.
		Insert(PullRequestTable).
		Columns(PullRequestColumns...).
		PlaceholderFormat(sq.Dollar).
		Values(pr.Values()...).
		ToSql()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("build insert pr sql: %v", err)
	}

	if _, err := tx.ExecContext(ctx, qPr, aPr...); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec insert pr: %v", err)
	}

	build := sq.
		Insert(ReviewersTable).
		Columns(ReviewerColumns...).
		PlaceholderFormat(sq.Dollar)

	for _, id := range reviewers {
		build = build.Values(model.Reviewer{
			PullRequestId: pr.Id,
			UserId:        id,
		}.Values()...)
	}

	qRev, aRev, err := build.ToSql()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("build reviewers insert sql: %v", err)
	}

	if _, err := tx.ExecContext(ctx, qRev, aRev...); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec reviewers insert: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}

	return nil
}

func (r *UserRepository) SetUserFlag(ctx context.Context, userID string, status bool) error {
	query, args, err := sq.
		Update(UsersTable).
		PlaceholderFormat(sq.Dollar).
		Set("is_active", status).
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update user sql: %v", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec update user: %v", err)
	}

	return nil
}

func (r *UserRepository) PullRequestReassign(ctx context.Context, rev model.Reviewer) error {
	qMeta, aMeta, err := sq.
		Select("users.team_id", "pull_requests.author_id").
		From(ReviewersTable).
		PlaceholderFormat(sq.Dollar).
		Join("pull_requests ON reviewers.pr_id = pull_requests.id").
		Join("users ON reviewers.user_id = users.id").
		Where(sq.And{
			sq.Eq{"pull_requests.status": "OPEN"},
			sq.Eq{"pull_requests.id": rev.PullRequestId},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build select for reassign: %v", err)
	}

	var teamID int64
	var authorID string

	if err := r.db.QueryRowContext(ctx, qMeta, aMeta...).Scan(&teamID, &authorID); err != nil {
		return fmt.Errorf("exec select for reassign: %v", err)
	}

	qNewRev, aNewRev, err := sq.
		Select("id").
		From(UsersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{"team_id": teamID},
			sq.NotEq{"id": rev.UserId},
			sq.NotEq{"id": authorID},
			sq.Eq{"is_active": true},
		}).
		OrderBy("random()").
		Limit(1).
		ToSql()
	if err != nil {
		return fmt.Errorf("build find new reviewer sql: %v", err)
	}

	var newReviewerID string
	if err := r.db.QueryRowContext(ctx, qNewRev, aNewRev...).Scan(&newReviewerID); err != nil {
		return fmt.Errorf("scan new reviewer: %v", err)
	}

	qUpdate, aUpdate, err := sq.
		Update(ReviewersTable).
		PlaceholderFormat(sq.Dollar).
		Set("pr_id", rev.PullRequestId).
		Set("user_id", newReviewerID).
		Where(sq.And{
			sq.Eq{"pr_id": rev.PullRequestId},
			sq.Eq{"user_id": rev.UserId},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build reviewer update sql: %v", err)
	}

	if _, err := r.db.ExecContext(ctx, qUpdate, aUpdate...); err != nil {
		return fmt.Errorf("exec reviewer update: %v", err)
	}

	return nil
}

func (r *UserRepository) MergePullRequest(ctx context.Context, prID string) error {
	query, args, err := sq.
		Update(PullRequestTable).
		PlaceholderFormat(sq.Dollar).
		Set("status", "MERGE").
		Where(sq.Eq{"id": prID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build merge sql: %v", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec merge: %v", err)
	}

	return nil
}

func (r *UserRepository) GetTeam(ctx context.Context, teamName string) ([]model.User, string, error) {
	query, args, err := sq.
		Select("users.id", "users.name", "is_active").
		PlaceholderFormat(sq.Dollar).
		From(UsersTable).
		Join("teams ON teams.id = users.team_id").
		Where(sq.Eq{"teams.name": teamName}).
		ToSql()
	if err != nil {
		return nil, "", fmt.Errorf("build get team sql: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("exec get team: %v", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var usr model.User
		if err := rows.Scan(&usr.Id, &usr.Name, &usr.IsActive); err != nil {
			return nil, "", fmt.Errorf("scan team: %v", err)
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("rows err: %v", err)
	}

	return users, teamName, nil
}

func (r *UserRepository) GetReview(ctx context.Context, userID string) ([]model.PullRequest, error) {
	query, args, err := sq.
		Select("id", "name", "author_id", "status").
		From(PullRequestTable).
		PlaceholderFormat(sq.Dollar).
		Join("reviewers ON reviewers.pr_id = pull_requests.id").
		Where(sq.Eq{"reviewers.user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get review sql: %v", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("exec get review: %v", err)
	}
	defer rows.Close()

	var res []model.PullRequest
	for rows.Next() {
		var p model.PullRequest
		if err := rows.Scan(&p.Id, &p.Name, &p.AuthorId, &p.Status); err != nil {
			return nil, fmt.Errorf("scan review: %v", err)
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %v", err)
	}

	return res, nil
}