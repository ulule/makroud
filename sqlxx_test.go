package sqlxx_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

var dbDefaultOptions = map[string]sqlxx.Option{
	"USER":     sqlxx.User("postgres"),
	"PASSWORD": sqlxx.Password(""),
	"HOST":     sqlxx.Host("localhost"),
	"PORT":     sqlxx.Port(5432),
	"NAME":     sqlxx.Database("sqlxx_test"),
}

type Elements struct {
	Air   string `db:"air"`
	Fire  string `sqlxx:"column:fire"`
	Water string `sqlxx:"-"`
	Earth string `sqlxx:"column:earth,default"`
	Fifth string
}

type Owl struct {
	ID           int64  `sqlxx:"column:id,pk"`
	Name         string `sqlxx:"column:name"`
	FeatherColor string `sqlxx:"column:feather_color"`
	FavoriteFood string
	tracking     bool
}

func (Owl) TableName() string {
	return "wp_owl"
}

//
// type Partner struct {
// 	ID   int
// 	Name string
// }
//
// func (Partner) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Partner", "partners").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name")
// }
//
// func (partner *Partner) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			partner.ID = int(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			partner.Name = value
// 		}),
// 	)
// }
//
// type PartnerList struct {
// 	partners []*Partner
// }
//
// func (e *PartnerList) Append(mapper sqlxx.Mapper) error {
// 	partner := &Partner{}
// 	err := partner.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if partner == nil {
// 		return errors.New("partner is nil")
// 	}
// 	e.partners = append(e.partners, partner)
// 	return nil
// }
//
// func (PartnerList) Model() sqlxx.XModel {
// 	return &Partner{}
// }
//
// func (e *PartnerList) One() *Partner {
// 	return e.partners[0]
// }
//
// func (e *PartnerList) List() []*Partner {
// 	return e.partners
// }
//
// type Manager struct {
// 	ID     int
// 	Name   string
// 	UserID int
// 	User   *User
// }
//
// func (Manager) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Manager", "managers").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name").
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
// 		AddAssociation("UserI", "User", sqlxx.AssociationTypeOne)
// }
//
// func (manager *Manager) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			manager.ID = int(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			manager.Name = value
// 		}),
// 		sqlxx.MapInt64("user_id", func(value int64) {
// 			manager.UserID = int(value)
// 		}),
// 	)
// }
//
// type ManagerList struct {
// 	managers []*Manager
// }
//
// func (e *ManagerList) Append(mapper sqlxx.Mapper) error {
// 	manager := &Manager{}
// 	err := manager.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if manager == nil {
// 		return errors.New("manager is nil")
// 	}
// 	e.managers = append(e.managers, manager)
// 	return nil
// }
//
// func (ManagerList) Model() sqlxx.XModel {
// 	return &Manager{}
// }
//
// func (e *ManagerList) One() *Manager {
// 	return e.managers[0]
// }
//
// func (e *ManagerList) List() []*Manager {
// 	return e.managers
// }
//
// type Project struct {
// 	ID        int
// 	Name      string
// 	ManagerID int
// 	UserID    int
// 	Manager   *Manager
// 	User      *User
// }
//
// func (Project) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Project", "projects").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name").
// 		AddField("ManagerID", "manager_id", sqlxx.IsForeignKey("Manager")).
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
// 		AddAssociation("Manager", "Manager", sqlxx.AssociationTypeOne).
// 		AddAssociation("User", "User", sqlxx.AssociationTypeOne)
// }
//
// func (project *Project) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			project.ID = int(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			project.Name = value
// 		}),
// 		sqlxx.MapInt64("user_id", func(value int64) {
// 			project.UserID = int(value)
// 		}),
// 		sqlxx.MapInt64("manager_id", func(value int64) {
// 			project.ManagerID = int(value)
// 		}),
// 	)
// }
//
// type ProjectList struct {
// 	projects []*Project
// }
//
// func (e *ProjectList) Append(mapper sqlxx.Mapper) error {
// 	project := &Project{}
// 	err := project.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if project == nil {
// 		return errors.New("project is nil")
// 	}
// 	e.projects = append(e.projects, project)
// 	return nil
// }
//
// func (ProjectList) Model() sqlxx.XModel {
// 	return &Project{}
// }
//
// func (e *ProjectList) One() *Project {
// 	return e.projects[0]
// }
//
// func (e *ProjectList) List() []*Project {
// 	return e.projects
// }
//
// type APIKey struct {
// 	ID        int
// 	Key       string
// 	Partner   Partner
// 	PartnerID int
// }
//
// func (APIKey) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("APIKey", "api_keys").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Key", "key").
// 		AddField("PartnerID", "partner_id", sqlxx.IsForeignKey("Partner")).
// 		AddAssociation("Partner", "Partner", sqlxx.AssociationTypeOne)
// }
//
// func (key *APIKey) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			key.ID = int(value)
// 		}),
// 		sqlxx.MapString("key", func(value string) {
// 			key.Key = value
// 		}),
// 		sqlxx.MapInt64("partner_id", func(value int64) {
// 			key.PartnerID = int(value)
// 		}),
// 	)
// }
//
// type APIKeyList struct {
// 	keys []*APIKey
// }
//
// func (e *APIKeyList) Append(mapper sqlxx.Mapper) error {
// 	key := &APIKey{}
// 	err := key.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if key == nil {
// 		return errors.New("key is nil")
// 	}
// 	e.keys = append(e.keys, key)
// 	return nil
// }
//
// func (APIKeyList) Model() sqlxx.XModel {
// 	return &APIKey{}
// }
//
// func (e *APIKeyList) One() *APIKey {
// 	return e.keys[0]
// }
//
// func (e *APIKeyList) List() []*APIKey {
// 	return e.keys
// }
//
// type Media struct {
// 	ID        int
// 	Path      string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
//
// func (Media) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Media", "media").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Path", "path").
// 		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()"))
// }
//
// func (media *Media) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			media.ID = int(value)
// 		}),
// 		sqlxx.MapString("path", func(value string) {
// 			media.Path = value
// 		}),
// 		sqlxx.MapTime("created_at", func(value time.Time) {
// 			media.CreatedAt = value
// 		}),
// 		sqlxx.MapTime("updated_at", func(value time.Time) {
// 			media.UpdatedAt = value
// 		}),
// 	)
// }
//
// type MediaList struct {
// 	medias []*Media
// }
//
// func (e *MediaList) Append(mapper sqlxx.Mapper) error {
// 	media := &Media{}
// 	err := media.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if media == nil {
// 		return errors.New("media is nil")
// 	}
// 	e.medias = append(e.medias, media)
// 	return nil
// }
//
// func (MediaList) Model() sqlxx.XModel {
// 	return &Media{}
// }
//
// func (e *MediaList) One() *Media {
// 	return e.medias[0]
// }
//
// func (e *MediaList) List() []*Media {
// 	return e.medias
// }
//
// type User struct {
// 	ID        int
// 	Username  string
// 	IsActive  bool
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time
//
// 	Avatars []Avatar
//
// 	Profile Profile
//
// 	APIKeyID int
// 	APIKey   APIKey
//
// 	NotificationID sql.NullInt64
// 	Notification   *Notification
//
// 	AvatarID sql.NullInt64
// 	Avatar   *Media
// }
//
// func (User) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("User", "users").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Username", "username").
// 		AddField("IsActive", "is_active", sqlxx.HasDefault("true")).
// 		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
// 		AddField("DeletedAt", "deleted_at", sqlxx.IsArchiveKey()).
// 		AddField("APIKeyID", "api_key_id", sqlxx.IsForeignKey("APIKey")).
// 		AddField("NotificationID", "notification_id", sqlxx.IsForeignKey("Notification")).
// 		AddField("AvatarID", "avatar_id", sqlxx.IsForeignKey("Media")).
// 		AddAssociation("Avatars", "Avatar", sqlxx.AssociationTypeMany).
// 		AddAssociation("Profile", "Profile", sqlxx.AssociationTypeOne).
// 		AddAssociation("APIKey", "APIKey", sqlxx.AssociationTypeOne).
// 		AddAssociation("Notification", "Notification", sqlxx.AssociationTypeOne).
// 		AddAssociation("Avatar", "Media", sqlxx.AssociationTypeOne)
// }
//
// func (user *User) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			user.ID = int(value)
// 		}),
// 		sqlxx.MapString("username", func(value string) {
// 			user.Username = value
// 		}),
// 		sqlxx.MapBool("is_active", func(value bool) {
// 			user.IsActive = value
// 		}),
// 		sqlxx.MapTime("created_at", func(value time.Time) {
// 			user.CreatedAt = value
// 		}),
// 		sqlxx.MapTime("updated_at", func(value time.Time) {
// 			user.UpdatedAt = value
// 		}),
// 		sqlxx.MapTime("deleted_at", func(value time.Time) {
// 			user.DeletedAt = &value
// 		}),
// 		sqlxx.MapInt64("api_key_id", func(value int64) {
// 			user.APIKeyID = int(value)
// 		}),
// 		sqlxx.MapNullInt64("notification_id", func(value sql.NullInt64) {
// 			user.NotificationID = value
// 		}),
// 		sqlxx.MapNullInt64("avatar_id", func(value sql.NullInt64) {
// 			user.AvatarID = value
// 		}),
// 	)
// }
//
// type UserList struct {
// 	users []*User
// }
//
// func (e *UserList) Append(mapper sqlxx.Mapper) error {
// 	user := &User{}
// 	err := user.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if user == nil {
// 		return errors.New("user is nil")
// 	}
// 	e.users = append(e.users, user)
// 	return nil
// }
//
// func (UserList) Model() sqlxx.XModel {
// 	return &User{}
// }
//
// func (e *UserList) One() *User {
// 	return e.users[0]
// }
//
// func (e *UserList) List() []*User {
// 	return e.users
// }
//
// type Notification struct {
// 	ID      int
// 	Enabled bool
// }
//
// func (Notification) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Notification", "notifications").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Enabled", "enabled", sqlxx.HasDefault("true"))
// }
//
// func (notification *Notification) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			notification.ID = int(value)
// 		}),
// 		sqlxx.MapBool("enabled", func(value bool) {
// 			notification.Enabled = value
// 		}),
// 	)
// }
//
// type Comment struct {
// 	ID        int
// 	UserID    int
// 	User      User
// 	ArticleID int
// 	Article   Article
// 	Content   string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
//
// func (Comment) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Comment", "comments").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Content", "content").
// 		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
// 		AddField("ArticleID", "article_id", sqlxx.IsForeignKey("Article")).
// 		AddAssociation("User", "User", sqlxx.AssociationTypeOne).
// 		AddAssociation("Article", "Article", sqlxx.AssociationTypeOne)
// }
//
// func (comment *Comment) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			comment.ID = int(value)
// 		}),
// 		sqlxx.MapInt64("user_id", func(value int64) {
// 			comment.UserID = int(value)
// 		}),
// 		sqlxx.MapInt64("article_id", func(value int64) {
// 			comment.ArticleID = int(value)
// 		}),
// 		sqlxx.MapString("content", func(value string) {
// 			comment.Content = value
// 		}),
// 		sqlxx.MapTime("created_at", func(value time.Time) {
// 			comment.CreatedAt = value
// 		}),
// 		sqlxx.MapTime("updated_at", func(value time.Time) {
// 			comment.UpdatedAt = value
// 		}),
// 	)
// }
//
// type Profile struct {
// 	ID        int64
// 	UserID    int
// 	FirstName string
// 	LastName  string
// }
//
// func (Profile) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Profile", "profiles").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("FirstName", "first_name").
// 		AddField("LastName", "last_name").
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User"))
// }
//
// func (profile *Profile) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			profile.ID = value
// 		}),
// 		sqlxx.MapInt64("user_id", func(value int64) {
// 			profile.UserID = int(value)
// 		}),
// 		sqlxx.MapString("first_name", func(value string) {
// 			profile.FirstName = value
// 		}),
// 		sqlxx.MapString("last_name", func(value string) {
// 			profile.LastName = value
// 		}),
// 	)
// }
//
// type ProfileList struct {
// 	profiles []*Profile
// }
//
// func (e *ProfileList) Append(mapper sqlxx.Mapper) error {
// 	profile := &Profile{}
// 	err := profile.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if profile == nil {
// 		return errors.New("profile is nil")
// 	}
// 	e.profiles = append(e.profiles, profile)
// 	return nil
// }
//
// func (ProfileList) Model() sqlxx.XModel {
// 	return &Profile{}
// }
//
// func (e *ProfileList) One() *Profile {
// 	return e.profiles[0]
// }
//
// func (e *ProfileList) List() []*Profile {
// 	return e.profiles
// }
//
// type AvatarFilter struct {
// 	ID   int
// 	Name string
// }
//
// func (AvatarFilter) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("AvatarFilter", "avatar_filters").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name")
// }
//
// func (filter *AvatarFilter) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			filter.ID = int(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			filter.Name = value
// 		}),
// 	)
// }
//
// type AvatarFilterList struct {
// 	filters []*AvatarFilter
// }
//
// func (e *AvatarFilterList) Append(mapper sqlxx.Mapper) error {
// 	filter := &AvatarFilter{}
// 	err := filter.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if filter == nil {
// 		return errors.New("filter is nil")
// 	}
// 	e.filters = append(e.filters, filter)
// 	return nil
// }
//
// func (AvatarFilterList) Model() sqlxx.XModel {
// 	return &AvatarFilter{}
// }
//
// func (e *AvatarFilterList) One() *AvatarFilter {
// 	return e.filters[0]
// }
//
// func (e *AvatarFilterList) List() []*AvatarFilter {
// 	return e.filters
// }
//
// type Avatar struct {
// 	ID        int
// 	Path      string
// 	UserID    int
// 	FilterID  int
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	Filter    *AvatarFilter
// }
//
// func (Avatar) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Avatar", "avatars").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Path", "path").
// 		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
// 		AddField("FilterID", "filter_id", sqlxx.IsForeignKey("AvatarFilter")).
// 		AddAssociation("Filter", "AvatarFilter", sqlxx.AssociationTypeOne)
// }
//
// func (avatar *Avatar) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			avatar.ID = int(value)
// 		}),
// 		sqlxx.MapString("path", func(value string) {
// 			avatar.Path = value
// 		}),
// 		sqlxx.MapInt64("user_id", func(value int64) {
// 			avatar.UserID = int(value)
// 		}),
// 		sqlxx.MapInt64("filter_id", func(value int64) {
// 			avatar.FilterID = int(value)
// 		}),
// 		sqlxx.MapTime("created_at", func(value time.Time) {
// 			avatar.CreatedAt = value
// 		}),
// 		sqlxx.MapTime("updated_at", func(value time.Time) {
// 			avatar.UpdatedAt = value
// 		}),
// 	)
// }
//
// type AvatarList struct {
// 	avatars []*Avatar
// }
//
// func (e *AvatarList) Append(mapper sqlxx.Mapper) error {
// 	avatar := &Avatar{}
// 	err := avatar.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if avatar == nil {
// 		return errors.New("avatar is nil")
// 	}
// 	e.avatars = append(e.avatars, avatar)
// 	return nil
// }
//
// func (AvatarList) Model() sqlxx.XModel {
// 	return &Avatar{}
// }
//
// func (e *AvatarList) One() *Avatar {
// 	return e.avatars[0]
// }
//
// func (e *AvatarList) List() []*Avatar {
// 	return e.avatars
// }
//
// type Category struct {
// 	ID     int           `db:"id" sqlxx:"primary_key:true"`
// 	Name   string        `db:"name"`
// 	UserID sql.NullInt64 `db:"user_id"`
// 	User   User
// }
//
// func (Category) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Category", "categories").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name").
// 		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
// 		AddAssociation("User", "User", sqlxx.AssociationTypeOne)
// }
//
// func (category *Category) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			category.ID = int(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			category.Name = value
// 		}),
// 		sqlxx.MapNullInt64("user_id", func(value sql.NullInt64) {
// 			category.UserID = value
// 		}),
// 	)
// }
//
// type CategoryList struct {
// 	categories []*Category
// }
//
// func (e *CategoryList) Append(mapper sqlxx.Mapper) error {
// 	category := &Category{}
// 	err := category.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if category == nil {
// 		return errors.New("category is nil")
// 	}
// 	e.categories = append(e.categories, category)
// 	return nil
// }
//
// func (CategoryList) Model() sqlxx.XModel {
// 	return &Category{}
// }
//
// func (e *CategoryList) One() *Category {
// 	return e.categories[0]
// }
//
// func (e *CategoryList) List() []*Category {
// 	return e.categories
// }
//
// // This model has a different ID type.
// type Tag struct {
// 	ID   uint
// 	Name string
// }
//
// func (Tag) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Tag", "tags").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Name", "name")
// }
//
// func (tag *Tag) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			tag.ID = uint(value)
// 		}),
// 		sqlxx.MapString("name", func(value string) {
// 			tag.Name = value
// 		}),
// 	)
// }
//
// type TagList struct {
// 	tags []*Tag
// }
//
// func (e *TagList) Append(mapper sqlxx.Mapper) error {
// 	tag := &Tag{}
// 	err := tag.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if tag == nil {
// 		return errors.New("tag is nil")
// 	}
// 	e.tags = append(e.tags, tag)
// 	return nil
// }
//
// func (TagList) Model() sqlxx.XModel {
// 	return &Tag{}
// }
//
// func (e *TagList) One() *Tag {
// 	return e.tags[0]
// }
//
// func (e *TagList) List() []*Tag {
// 	return e.tags
// }
//
// type Article struct {
// 	ID          int
// 	Title       string
// 	AuthorID    int
// 	ReviewerID  int
// 	IsPublished bool
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// 	Author      User
// 	Reviewer    *User
// 	MainTagID   sql.NullInt64
// 	MainTag     *Tag
// }
//
// func (Article) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("Article", "articles").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("Title", "title").
// 		AddField("IsPublished", "is_published").
// 		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
// 		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
// 		AddField("AuthorID", "author_id", sqlxx.IsForeignKey("User")).
// 		AddField("ReviewerID", "reviewer_id", sqlxx.IsForeignKey("User")).
// 		AddField("MainTagID", "main_tag_id", sqlxx.IsForeignKey("Tag")).
// 		AddAssociation("Author", "User", sqlxx.AssociationTypeOne).
// 		AddAssociation("Reviewer", "User", sqlxx.AssociationTypeOne).
// 		AddAssociation("MainTag", "Tag", sqlxx.AssociationTypeOne)
// }
//
// func (article *Article) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			article.ID = int(value)
// 		}),
// 		sqlxx.MapString("title", func(value string) {
// 			article.Title = value
// 		}),
// 		sqlxx.MapBool("is_published", func(value bool) {
// 			article.IsPublished = value
// 		}),
// 		sqlxx.MapInt64("author_id", func(value int64) {
// 			article.AuthorID = int(value)
// 		}),
// 		sqlxx.MapInt64("reviewer_id", func(value int64) {
// 			article.ReviewerID = int(value)
// 		}),
// 		sqlxx.MapNullInt64("main_tag_id", func(value sql.NullInt64) {
// 			article.MainTagID = value
// 		}),
// 		sqlxx.MapTime("created_at", func(value time.Time) {
// 			article.CreatedAt = value
// 		}),
// 		sqlxx.MapTime("updated_at", func(value time.Time) {
// 			article.UpdatedAt = value
// 		}),
// 	)
// }
//
// type ArticleList struct {
// 	articles []*Article
// }
//
// func (e *ArticleList) Append(mapper sqlxx.Mapper) error {
// 	article := &Article{}
// 	err := article.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if article == nil {
// 		return errors.New("article is nil")
// 	}
// 	e.articles = append(e.articles, article)
// 	return nil
// }
//
// func (ArticleList) Model() sqlxx.XModel {
// 	return &Article{}
// }
//
// func (e *ArticleList) One() *Article {
// 	return e.articles[0]
// }
//
// func (e *ArticleList) List() []*Article {
// 	return e.articles
// }
//
// type ArticleCategory struct {
// 	ID         int
// 	ArticleID  int
// 	CategoryID int
// }
//
// func (ArticleCategory) CreateSchema(builder sqlxx.SchemaBuilder) {
// 	builder.SetTableName("ArticleCategory", "articles_categories").
// 		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
// 		AddField("ArticleID", "article_id", sqlxx.IsForeignKey("Article")).
// 		AddField("CategoryID", "category_id", sqlxx.IsForeignKey("Category"))
// }
//
// func (article *ArticleCategory) WriteModel(mapper sqlxx.Mapper) error {
// 	return sqlxx.Map(mapper,
// 		sqlxx.MapInt64("id", func(value int64) {
// 			article.ID = int(value)
// 		}),
// 		sqlxx.MapInt64("article_id", func(value int64) {
// 			article.ArticleID = int(value)
// 		}),
// 		sqlxx.MapInt64("category_id", func(value int64) {
// 			article.CategoryID = int(value)
// 		}),
// 	)
// }
//
// type ArticleCategoryList struct {
// 	relations []*ArticleCategory
// }
//
// func (e *ArticleCategoryList) Append(mapper sqlxx.Mapper) error {
// 	relation := &ArticleCategory{}
// 	err := relation.WriteModel(mapper)
// 	if err != nil {
// 		return err
// 	}
// 	if relation == nil {
// 		return errors.New("article-category is nil")
// 	}
// 	e.relations = append(e.relations, relation)
// 	return nil
// }
//
// func (ArticleCategoryList) Model() sqlxx.XModel {
// 	return &ArticleCategory{}
// }
//
// func (e *ArticleCategoryList) One() *ArticleCategory {
// 	return e.relations[0]
// }
//
// func (e *ArticleCategoryList) List() []*ArticleCategory {
// 	return e.relations
// }

// ----------------------------------------------------------------------------
// Loader
// ----------------------------------------------------------------------------

type environment struct {
	driver *sqlxx.Client
	is     *require.Assertions
	// Users              *UserList
	// APIKeys            *APIKeyList
	// Profiles           *ProfileList
	// AvatarFilters      *AvatarFilterList
	// Avatars            *AvatarList
	// Articles           *ArticleList
	// Categories         *CategoryList
	// Tags               *TagList
	// ArticlesCategories *ArticleCategoryList
	// Partners           *PartnerList
	// Managers           *ManagerList
	// Projects           *ProjectList
	// Medias             *MediaList
}

func (e *environment) startup() {
	DropTables(e.driver)
	CreateTables(e.driver)

	// e.insertPartners()
	// e.insertAPIKeys()
	// e.insertMedias()
	// e.insertUsers()
	// e.insertManagers()
	// e.insertProjects()
	// e.insertAvatarFilters()
	// e.insertAvatars()
	// e.insertProfiles()
	// e.insertCategories()
	// e.insertTags()
	// e.insertArticles()
}

func (e *environment) fetch(query string, callback func(mapper sqlxx.Mapper)) {
	is := e.is
	driver := e.driver

	rows, err := driver.Queryx(query)
	is.NoError(err)
	is.NotNil(rows)
	defer rows.Close()

	for rows.Next() {
		mapper := map[string]interface{}{}
		err = rows.MapScan(mapper)
		is.NoError(err)
		callback(mapper)
	}
	err = rows.Err()
	is.NoError(err)
}

func (e *environment) exec(query string, args ...interface{}) {
	e.driver.MustExec(query, args...)
}

// func (e *environment) insertPartners() {
// 	list := &PartnerList{}
//
// 	e.exec(`INSERT INTO partners (name) VALUES ($1)`,
// 		"Wayne Enterprise",
// 	)
//
// 	e.fetch(`SELECT * FROM partners`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Partners = list
// }
//
// func (e *environment) insertAPIKeys() {
// 	list := &APIKeyList{}
//
// 	e.exec(`INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)`,
// 		"caf57d0bcaae636f0a9a5d072da6f096990b37606d724a5b65e43e153f5eddb8",
// 		e.Partners.One().ID,
// 	)
//
// 	e.fetch(`SELECT * FROM api_keys`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.APIKeys = list
// }
//
// func (e *environment) insertMedias() {
// 	list := &MediaList{}
//
// 	e.exec(`INSERT INTO media (path) VALUES ($1)`,
// 		"media/avatar.png",
// 	)
//
// 	e.fetch(`SELECT * FROM media`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Medias = list
// }
// func (e *environment) insertUsers() {
// 	list := &UserList{}
//
// 	e.exec(`INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)`,
// 		"lucius_fox",
// 		e.APIKeys.One().ID,
// 		e.Medias.One().ID,
// 	)
//
// 	e.fetch(`SELECT * FROM users`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Users = list
// }
//
// func (e *environment) insertManagers() {
// 	list := &ManagerList{}
//
// 	e.exec(`INSERT INTO managers (name, user_id) VALUES ($1, $2)`,
// 		"Lucius Fox",
// 		e.Users.One().ID,
// 	)
//
// 	e.fetch(`SELECT * FROM managers`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Managers = list
// }
//
// func (e *environment) insertProjects() {
// 	list := &ProjectList{}
//
// 	e.exec(`INSERT INTO projects (name, manager_id, user_id) VALUES ($1, $2, $3)`,
// 		"Super Project",
// 		e.Managers.One().ID,
// 		e.Users.One().ID,
// 	)
//
// 	e.fetch(`SELECT * FROM projects`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Projects = list
// }
//
// func (e *environment) insertAvatarFilters() {
// 	list := &AvatarFilterList{}
//
// 	names := []string{"normal", "clarendon", "juno", "lark", "ludwig", "gingham", "valencia"}
// 	for _, name := range names {
// 		e.exec(`INSERT INTO avatar_filters (name) VALUES ($1)`, name)
// 	}
//
// 	e.fetch(`SELECT * FROM avatar_filters`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.AvatarFilters = list
// }
//
// func (e *environment) insertAvatars() {
// 	list := &AvatarList{}
//
// 	for i := 0; i < 5; i++ {
// 		e.exec(`INSERT INTO avatars (path, user_id, filter_id) VALUES ($1, $2, $3)`,
// 			fmt.Sprintf("/avatars/%s-%d.png", e.Users.One().Username, i),
// 			e.Users.One().ID,
// 			e.AvatarFilters.One().ID,
// 		)
// 	}
// 	e.fetch(`SELECT * FROM avatars`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Avatars = list
// }
//
// func (e *environment) insertProfiles() {
// 	list := &ProfileList{}
//
// 	e.exec(`INSERT INTO profiles (user_id, first_name, last_name) VALUES ($1, $2, $3)`,
// 		e.Users.One().ID,
// 		"Lucius",
// 		"Fox",
// 	)
//
// 	e.fetch(`SELECT * FROM profiles`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Profiles = list
// }
//
// func (e *environment) insertCategories() {
// 	list := &CategoryList{}
//
// 	for i := 0; i < 5; i++ {
// 		e.exec(`INSERT INTO categories (name, user_id) VALUES ($1, $2)`,
// 			fmt.Sprintf("Category #%d", i),
// 			e.Users.One().ID,
// 		)
// 	}
//
// 	e.fetch(`SELECT * FROM categories`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Categories = list
// }
//
// func (e *environment) insertTags() {
// 	list := &TagList{}
//
// 	e.exec(`INSERT INTO tags (name) VALUES ($1)`, "Tag")
//
// 	e.fetch(`SELECT * FROM tags`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Tags = list
// }
//
// func (e *environment) insertArticles() {
// 	list := &ArticleList{}
// 	relations := &ArticleCategoryList{}
//
// 	for i := 0; i < 5; i++ {
// 		e.exec(`INSERT INTO articles (title, author_id, reviewer_id, main_tag_id) VALUES ($1, $2, $3, $4)`,
// 			fmt.Sprintf("Title #%d", i),
// 			e.Users.One().ID,
// 			e.Users.One().ID,
// 			e.Tags.One().ID,
// 		)
// 	}
//
// 	e.fetch(`SELECT * FROM articles`, func(mapper sqlxx.Mapper) {
// 		err := list.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.Articles = list
//
// 	for _, article := range e.Articles.List() {
// 		for _, category := range e.Categories.List() {
// 			e.exec(`INSERT INTO articles_categories (article_id, category_id) VALUES ($1, $2)`,
// 				article.ID,
// 				category.ID,
// 			)
// 		}
// 	}
//
// 	e.fetch(`SELECT * FROM articles_categories`, func(mapper sqlxx.Mapper) {
// 		err := relations.Append(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	e.ArticlesCategories = relations
// }
//
// func (e *environment) createComment(user *User, article *Article) *Comment {
// 	comment := &Comment{}
// 	e.exec(`INSERT INTO comments (content, user_id, article_id) VALUES ($1, $2, $3);`,
// 		"Lorem Ipsum",
// 		user.ID,
// 		article.ID,
// 	)
// 	e.fetch(`SELECT * FROM comments ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := comment.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
// 	return comment
// }
//
// func (e *environment) createArticle(user *User) *Article {
// 	article := &Article{}
// 	e.exec(`INSERT INTO articles (title, author_id, reviewer_id) VALUES ($1, $2, $3);`,
// 		"Title",
// 		user.ID,
// 		user.ID,
// 	)
// 	e.fetch(`SELECT * FROM articles ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := article.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	return article
// }
//
// func (e *environment) createUser(username string) *User {
// 	key := fmt.Sprintf("%s-apikey", username)
// 	name := fmt.Sprintf("%s-partner", username)
// 	path := fmt.Sprintf("media/media-%s.png", username)
//
// 	partner := &Partner{}
// 	e.exec(`INSERT INTO partners (name) VALUES ($1)`, name)
// 	e.fetch(`SELECT * FROM partners ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := partner.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	media := &Media{}
// 	e.exec(`INSERT INTO media (path) VALUES ($1)`, path)
// 	e.fetch(`SELECT * FROM media ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := media.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	apiKey := &APIKey{}
// 	e.exec(`INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)`, key, partner.ID)
// 	e.fetch(`SELECT * FROM api_keys ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := apiKey.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	user := &User{}
// 	e.exec(`INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)`,
// 		username,
// 		apiKey.ID,
// 		media.ID,
// 	)
// 	e.fetch(`SELECT * FROM users ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := user.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	for filterID := 1; filterID < 6; filterID++ {
// 		e.exec(`INSERT INTO avatars (path, user_id, filter_id) VALUES ($1, $2, $3)`,
// 			fmt.Sprintf("/avatars/%s-%d.png", username, filterID),
// 			user.ID,
// 			filterID,
// 		)
// 	}
//
// 	return user
// }
//
// func (e *environment) createCategory(name string, userID *int) *Category {
// 	e.exec(`INSERT INTO categories (name) VALUES ($1)`, name)
// 	if userID != nil {
// 		e.exec(`UPDATE categories SET user_id=$1 WHERE name=$2`, *userID, name)
// 	}
//
// 	category := &Category{}
// 	e.fetch(`SELECT * FROM categories ORDER BY id DESC LIMIT 1`, func(mapper sqlxx.Mapper) {
// 		err := category.WriteModel(mapper)
// 		e.is.NoError(err)
// 	})
//
// 	return category
// }

func (e *environment) shutdown() {
	value := os.Getenv("KEEP_DB")
	if len(value) == 0 {
		DropTables(e.driver)
	}
	e.is.NoError(e.driver.Close())
}

func dbParamString(option func(string) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	if len(v) != 0 {
		return option(v)
	}
	for i := range env {
		v = os.Getenv(env[i])
		if len(v) != 0 {
			return option(v)
		}
	}
	return dbDefaultOptions[param]
}

func dbParamInt(option func(int) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	n, err := strconv.Atoi(v)
	if err == nil {
		return option(n)
	}
	for i := range env {
		v = os.Getenv(env[i])
		n, err = strconv.Atoi(v)
		if err == nil {
			return option(n)
		}
	}
	return dbDefaultOptions[param]
}

type SetupCallback func(handler SetupHandler)

type SetupHandler func(driver sqlxx.Driver)

func Setup(t *testing.T, options ...sqlxx.Option) SetupCallback {
	is := require.New(t)
	opts := []sqlxx.Option{
		dbParamString(sqlxx.Host, "host", "PGHOST"),
		dbParamInt(sqlxx.Port, "port", "PGPORT"),
		dbParamString(sqlxx.User, "user", "PGUSER"),
		dbParamString(sqlxx.Password, "password", "PGPASSWORD"),
		dbParamString(sqlxx.Database, "name", "PGDATABASE"),
		sqlxx.Cache(false),
	}
	opts = append(opts, options...)

	db, err := sqlxx.New(opts...)
	is.NoError(err)
	is.NotNil(db)

	env := &environment{
		is:     is,
		driver: db,
	}

	return func(handler SetupHandler) {
		env.startup()
		handler(db)
		env.shutdown()
	}
}

func DropTables(db *sqlxx.Client) {
	db.MustExec(`
		-- Simple schema
		DROP TABLE IF EXISTS wp_owl CASCADE;

		-- Application schema
		DROP TABLE IF EXISTS users CASCADE;
		DROP TABLE IF EXISTS api_keys CASCADE;
		DROP TABLE IF EXISTS profiles CASCADE;
		DROP TABLE IF EXISTS comments CASCADE;
		DROP TABLE IF EXISTS avatar_filters CASCADE;
		DROP TABLE IF EXISTS avatars CASCADE;
		DROP TABLE IF EXISTS categories CASCADE;
		DROP TABLE IF EXISTS tags CASCADE;
		DROP TABLE IF EXISTS articles CASCADE;
		DROP TABLE IF EXISTS articles_categories CASCADE;
		DROP TABLE IF EXISTS partners CASCADE;
		DROP TABLE IF EXISTS media CASCADE;
		DROP TABLE IF EXISTS projects CASCADE;
		DROP TABLE IF EXISTS managers CASCADE;
		DROP TABLE IF EXISTS notifications CASCADE;
	`)
}

func CreateTables(db *sqlxx.Client) {
	db.MustExec(`
		-- Simple schema
		CREATE TABLE wp_owl (
			id              SERIAL PRIMARY KEY NOT NULL,
			name            VARCHAR(255) NOT NULL,
			feather_color   VARCHAR(255) NOT NULL,
			favorite_food   VARCHAR(255) NOT NULL
		);


		-- Application schema
		CREATE TABLE notifications (
			id       serial primary key not null,
			enabled  boolean default true
		);

		CREATE TABLE api_keys (
			id          serial primary key not null,
			partner_id  integer,
			key         varchar(255) not null
		);

		CREATE TABLE partners (
			id    serial primary key not null,
			name  varchar(255) not null
		);

		CREATE TABLE managers (
			id       serial primary key not null,
			name     varchar(255) not null,
			user_id  integer
		);

		CREATE TABLE projects (
			id          serial primary key not null,
			name        varchar(255) not null,
			manager_id  integer,
			user_id     integer
		);

		CREATE TABLE users (
			id               serial primary key not null,
			username         varchar(30) not null,
			is_active        boolean default true,
			api_key_id       integer,
			avatar_id        integer,
			notification_id  integer references notifications(id) default null,
			created_at       timestamp with time zone default current_timestamp,
			updated_at       timestamp with time zone default current_timestamp,
			deleted_at       timestamp with time zone
		);

		CREATE TABLE profiles (
			id          serial primary key not null,
			user_id     integer references users(id),
			first_name  varchar(255) not null,
			last_name   varchar(255) not null
		);

		CREATE TABLE media (
			id          serial primary key not null,
			path        varchar(255) not null,
			created_at  timestamp with time zone default current_timestamp,
			updated_at  timestamp with time zone default current_timestamp
		);

		CREATE TABLE tags (
			id    serial primary key not null,
			name  varchar(255) not null
		);

		CREATE TABLE avatar_filters (
			id    serial primary key not null,
			name  varchar(255) not null
		);

		CREATE TABLE avatars (
			id          serial primary key not null,
			path        varchar(255) not null,
			user_id     integer references users(id),
			filter_id   integer references avatar_filters(id),
			created_at  timestamp with time zone default current_timestamp,
			updated_at  timestamp with time zone default current_timestamp
		);

		CREATE TABLE articles (
			id            serial primary key not null,
			title         varchar(255) not null,
			author_id     integer references users(id),
			reviewer_id   integer references users(id),
			main_tag_id   integer references tags(id),
			is_published  boolean default true,
			created_at    timestamp with time zone default current_timestamp,
			updated_at    timestamp with time zone default current_timestamp
		);

		CREATE TABLE comments (
			id          serial primary key not null,
			user_id     integer references users(id),
			article_id  integer references articles(id),
			content     text,
			created_at  timestamp with time zone default current_timestamp,
			updated_at  timestamp with time zone default current_timestamp
		);

		CREATE TABLE categories (
			id       serial primary key not null,
			name     varchar(255) not null,
			user_id  integer references users(id)
		);

		CREATE TABLE articles_categories (
			id           serial primary key not null,
			article_id   integer references articles(id),
			category_id  integer references categories(id)
		);
	`)
}
