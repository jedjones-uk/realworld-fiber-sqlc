package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

type MockQuerier struct{}

func (m *MockQuerier) CreateArticle(ctx context.Context, arg *CreateArticleParams) (CreateArticleRow, error) {
	return CreateArticleRow{}, nil
}

func (m *MockQuerier) CreateComment(ctx context.Context, arg *CreateCommentParams) (CreateCommentRow, error) {
	return CreateCommentRow{}, nil
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg *CreateUserParams) (User, error) {
	return User{}, nil
}

func (m *MockQuerier) DeleteArticle(ctx context.Context, arg *DeleteArticleParams) error {
	return nil
}

func (m *MockQuerier) DeleteComment(ctx context.Context, arg *DeleteCommentParams) error {
	return nil
}

func (m *MockQuerier) FavoriteArticle(ctx context.Context, arg *FavoriteArticleParams) (FavoriteArticleRow, error) {
	return FavoriteArticleRow{}, nil
}

func (m *MockQuerier) FeedArticles(ctx context.Context, arg *FeedArticlesParams) ([]FeedArticlesRow, error) {
	return nil, nil
}

func (m *MockQuerier) FollowUser(ctx context.Context, arg *FollowUserParams) (FollowUserRow, error) {
	return FollowUserRow{}, nil
}

func (m *MockQuerier) GetArticle(ctx context.Context, arg *GetArticleParams) (GetArticleRow, error) {
	return GetArticleRow{}, nil
}

func (m *MockQuerier) GetCommentsByArticleSlug(ctx context.Context, slug string) ([]GetCommentsByArticleSlugRow, error) {
	return nil, nil
}

func (m *MockQuerier) GetSingleComment(ctx context.Context) (GetSingleCommentRow, error) {
	return GetSingleCommentRow{}, nil
}

func (m *MockQuerier) GetTags(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (m *MockQuerier) GetUser(ctx context.Context, id int64) (GetUserRow, error) {
	return GetUserRow{
		Email:    "test@example.com",
		Username: "testuser",
		Bio:      pgtype.Text{String: "test bio", Valid: true},
		Image:    pgtype.Text{String: "test image", Valid: true},
	}, nil
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return User{}, nil
}

func (m *MockQuerier) GetUserProfile(ctx context.Context, arg *GetUserProfileParams) (GetUserProfileRow, error) {
	return GetUserProfileRow{}, nil
}

func (m *MockQuerier) GetUserProfileById(ctx context.Context, arg *GetUserProfileByIdParams) (GetUserProfileByIdRow, error) {
	return GetUserProfileByIdRow{}, nil
}

func (m *MockQuerier) ListArticles(ctx context.Context, arg *ListArticlesParams) ([]ListArticlesRow, error) {
	return nil, nil
}

func (m *MockQuerier) UnfavoriteArticle(ctx context.Context, arg *UnfavoriteArticleParams) (UnfavoriteArticleRow, error) {
	return UnfavoriteArticleRow{}, nil
}

func (m *MockQuerier) UnfollowUser(ctx context.Context, arg *UnfollowUserParams) (UnfollowUserRow, error) {
	return UnfollowUserRow{}, nil
}

func (m *MockQuerier) UpdateArticle(ctx context.Context, arg *UpdateArticleParams) (UpdateArticleRow, error) {
	return UpdateArticleRow{}, nil
}

func (m *MockQuerier) UpdateUser(ctx context.Context, arg *UpdateUserParams) (UpdateUserRow, error) {
	return UpdateUserRow{}, nil
}
