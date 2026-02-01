package models

import (
	"math/rand/v2"
	"time"
	"fmt"
	"github.com/google/uuid"
)

var shortWords = []string{
	"big", "new", "old", "hot", "red", "top", "far", "fun", "sky", "sun",
	"bold", "calm", "cool", "dark", "deep", "easy", "fair", "fast", "fine",
	"free", "full", "good", "gray", "hard", "high", "kind", "last", "lazy",
	"long", "loud", "near", "next", "nice", "pink", "pure", "rare", "real",
	"rich", "safe", "slow", "soft", "tall", "thin", "tidy", "tiny", "true",
	"ugly", "vast", "warm", "weak", "wild", "wise", "zero",
}

var mediumWords = []string{
	"amber", "beach", "brave", "bread", "bring", "chair", "charm", "chess",
	"clear", "crown", "dance", "dream", "eagle", "flame", "frost", "glass",
	"grace", "grape", "green", "happy", "heart", "honey", "horse", "house",
	"light", "maple", "merit", "metal", "music", "noble", "ocean", "olive",
	"peace", "pearl", "piano", "plant", "queen", "quick", "quiet", "river",
	"royal", "smart", "snake", "space", "spark", "storm", "sweet", "swift",
	"tiger", "trust", "unity", "urban", "value", "water", "whale", "world",
}


// Account represents the main account/organization
type Account struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	Status    string    `json:"status" db:"status" gorm:"type:varchar(50);default:'active'" validate:"oneof=active suspended deleted"`
	
	// Relationships
	Tenants      []Tenant      `json:"tenants,omitempty" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	AccountUsers []AccountUser `json:"account_users,omitempty" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
}

// Tenant represents an isolated environment within an account
type Tenant struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AccountID uuid.UUID `json:"account_id" db:"account_id" gorm:"type:uuid;not null;index" validate:"required"`
	Slug      string    `json:"slug" db:"slug" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	Status    string    `json:"status" db:"status" gorm:"type:varchar(50);default:'active'" validate:"oneof=active suspended deleted"`
	Settings  []byte    `json:"settings,omitempty" db:"settings" gorm:"type:jsonb"` // Use json.RawMessage or custom type
	
	// Relationships
	Account    Account     `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Clients    []Client    `json:"clients,omitempty" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Users      []User      `json:"users,omitempty" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	AuditLogs  []AuditLog  `json:"audit_logs,omitempty" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
}

// Client represents an OAuth client application
type Client struct {
	ID               uuid.UUID `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClientID         string    `json:"client_id" db:"client_id" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	ClientSecretHash string    `json:"-" db:"client_secret_hash" gorm:"type:varchar(255);not null" validate:"required"` // Never expose in JSON
	TenantID         uuid.UUID `json:"tenant_id" db:"tenant_id" gorm:"type:uuid;not null;index" validate:"required"`
	Name             string    `json:"name" db:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Description      string    `json:"description,omitempty" db:"description" gorm:"type:text"`
	RedirectURIs     string    `json:"redirect_uris" db:"redirect_uris" gorm:"type:text;not null" validate:"required"` // Store as JSON string
	GrantTypes       string    `json:"grant_types" db:"grant_types" gorm:"type:text;not null" validate:"required"`
	ResponseTypes    string    `json:"response_types" db:"response_types" gorm:"type:text;not null" validate:"required"`
	Scopes           string    `json:"scopes" db:"scopes" gorm:"type:text;not null" validate:"required"`
	IsConfidential   bool      `json:"is_confidential" db:"is_confidential" gorm:"default:true"`
	CreatedAt        time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	Status           string    `json:"status" db:"status" gorm:"type:varchar(50);default:'active'" validate:"oneof=active suspended deleted"`
	
	// Relationships
	Tenant              Tenant               `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	AuthorizationCodes  []AuthorizationCode  `json:"authorization_codes,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	AccessTokens        []AccessToken        `json:"access_tokens,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	RefreshTokens       []RefreshToken       `json:"refresh_tokens,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	IDTokens            []IDToken            `json:"id_tokens,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	Sessions            []Session            `json:"sessions,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	UserConsents        []UserConsent        `json:"user_consents,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

// User represents an end user with OpenID Connect identity
type User struct {
	ID            uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID      uuid.UUID  `json:"tenant_id" db:"tenant_id" gorm:"type:uuid;not null;index" validate:"required"`
	Email         string     `json:"email" db:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	EmailVerified bool       `json:"email_verified" db:"email_verified" gorm:"default:false"`
	PasswordHash  string     `json:"-" db:"password_hash" gorm:"type:varchar(255)"` // Null for social logins
	GivenName     string     `json:"given_name,omitempty" db:"given_name" gorm:"type:varchar(255)"`
	FamilyName    string     `json:"family_name,omitempty" db:"family_name" gorm:"type:varchar(255)"`
	PictureURL    string     `json:"picture_url,omitempty" db:"picture_url" gorm:"type:text"`
	Locale        string     `json:"locale,omitempty" db:"locale" gorm:"type:varchar(10)"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	Status        string     `json:"status" db:"status" gorm:"type:varchar(50);default:'active'" validate:"oneof=active suspended deleted"`
	
	// Relationships
	Tenant             Tenant              `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	AuthorizationCodes []AuthorizationCode `json:"authorization_codes,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AccessTokens       []AccessToken       `json:"access_tokens,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RefreshTokens      []RefreshToken      `json:"refresh_tokens,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	IDTokens           []IDToken           `json:"id_tokens,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Sessions           []Session           `json:"sessions,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	UserConsents       []UserConsent       `json:"user_consents,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AuditLogs          []AuditLog          `json:"audit_logs,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

// AuthorizationCode represents a short-lived authorization code
type AuthorizationCode struct {
	ID                    uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Code                  string     `json:"code" db:"code" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	ClientID              uuid.UUID  `json:"client_id" db:"client_id" gorm:"type:uuid;not null;index" validate:"required"`
	UserID                uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	RedirectURI           string     `json:"redirect_uri" db:"redirect_uri" gorm:"type:text;not null" validate:"required,url"`
	Scopes                string     `json:"scopes" db:"scopes" gorm:"type:text;not null" validate:"required"`
	CodeChallenge         string     `json:"code_challenge,omitempty" db:"code_challenge" gorm:"type:varchar(255)"`
	CodeChallengeMethod   string     `json:"code_challenge_method,omitempty" db:"code_challenge_method" gorm:"type:varchar(10)" validate:"omitempty,oneof=S256 plain"`
	ExpiresAt             time.Time  `json:"expires_at" db:"expires_at" gorm:"not null;index" validate:"required"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UsedAt                *time.Time `json:"used_at,omitempty" db:"used_at"`
	Nonce                 string     `json:"nonce,omitempty" db:"nonce" gorm:"type:varchar(255)"`
	
	// Relationships
	Client   Client    `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	IDTokens []IDToken `json:"id_tokens,omitempty" gorm:"foreignKey:AuthorizationCodeID"`
}

// AccessToken represents an OAuth access token
type AccessToken struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash string     `json:"-" db:"token_hash" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	ClientID  uuid.UUID  `json:"client_id" db:"client_id" gorm:"type:uuid;not null;index" validate:"required"`
	UserID    *uuid.UUID `json:"user_id,omitempty" db:"user_id" gorm:"type:uuid;index"` // Null for client_credentials
	Scopes    string     `json:"scopes" db:"scopes" gorm:"type:text;not null" validate:"required"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at" gorm:"not null;index" validate:"required"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	SessionID *uuid.UUID `json:"session_id,omitempty" db:"session_id" gorm:"type:uuid;index"`
	
	// Relationships
	Client        Client         `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Session       *Session       `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	RefreshTokens []RefreshToken `json:"refresh_tokens,omitempty" gorm:"foreignKey:AccessTokenID"`
}

// RefreshToken represents an OAuth refresh token
type RefreshToken struct {
	ID            uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash     string     `json:"-" db:"token_hash" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	AccessTokenID *uuid.UUID `json:"access_token_id,omitempty" db:"access_token_id" gorm:"type:uuid"`
	ClientID      uuid.UUID  `json:"client_id" db:"client_id" gorm:"type:uuid;not null;index" validate:"required"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	Scopes        string     `json:"scopes" db:"scopes" gorm:"type:text;not null" validate:"required"`
	ExpiresAt     time.Time  `json:"expires_at" db:"expires_at" gorm:"not null" validate:"required"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	SessionID     *uuid.UUID `json:"session_id,omitempty" db:"session_id" gorm:"type:uuid;index"`
	
	// Relationships
	AccessToken *AccessToken `json:"access_token,omitempty" gorm:"foreignKey:AccessTokenID"`
	Client      Client       `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	User        User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Session     *Session     `json:"session,omitempty" gorm:"foreignKey:SessionID"`
}

// IDToken represents an OpenID Connect ID token
type IDToken struct {
	ID                   uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash            string     `json:"-" db:"token_hash" gorm:"type:varchar(255);not null;uniqueIndex" validate:"required"`
	AuthorizationCodeID  *uuid.UUID `json:"authorization_code_id,omitempty" db:"authorization_code_id" gorm:"type:uuid;index"`
	ClientID             uuid.UUID  `json:"client_id" db:"client_id" gorm:"type:uuid;not null;index" validate:"required"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	Nonce                string     `json:"nonce,omitempty" db:"nonce" gorm:"type:varchar(255)"`
	ACR                  string     `json:"acr,omitempty" db:"acr" gorm:"type:varchar(50)"`
	AMR                  string     `json:"amr,omitempty" db:"amr" gorm:"type:text"`
	AZP                  string     `json:"azp,omitempty" db:"azp" gorm:"type:varchar(255)"`
	ATHash               string     `json:"at_hash,omitempty" db:"at_hash" gorm:"type:varchar(255)"`
	CHash                string     `json:"c_hash,omitempty" db:"c_hash" gorm:"type:varchar(255)"`
	IssuedAt             time.Time  `json:"issued_at" db:"issued_at" gorm:"not null" validate:"required"`
	ExpiresAt            time.Time  `json:"expires_at" db:"expires_at" gorm:"not null;index" validate:"required"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	SessionID            *uuid.UUID `json:"session_id,omitempty" db:"session_id" gorm:"type:uuid;index"`
	
	// Relationships
	AuthorizationCode *AuthorizationCode `json:"authorization_code,omitempty" gorm:"foreignKey:AuthorizationCodeID"`
	Client            Client             `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	User              User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Session           *Session           `json:"session,omitempty" gorm:"foreignKey:SessionID"`
}

// Session represents a user session
type Session struct {
	ID             uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	ClientID       *uuid.UUID `json:"client_id,omitempty" db:"client_id" gorm:"type:uuid;index"`
	UserAgent      string     `json:"user_agent,omitempty" db:"user_agent" gorm:"type:text"`
	IPAddress      string     `json:"ip_address,omitempty" db:"ip_address" gorm:"type:varchar(45)"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	LastActivityAt time.Time  `json:"last_activity_at" db:"last_activity_at" gorm:"autoCreateTime"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at" gorm:"not null;index" validate:"required"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	
	// Relationships
	User          User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Client        *Client        `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	AccessTokens  []AccessToken  `json:"access_tokens,omitempty" gorm:"foreignKey:SessionID"`
	RefreshTokens []RefreshToken `json:"refresh_tokens,omitempty" gorm:"foreignKey:SessionID"`
	IDTokens      []IDToken      `json:"id_tokens,omitempty" gorm:"foreignKey:SessionID"`
}

// UserConsent represents user consent to client access
type UserConsent struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	ClientID  uuid.UUID  `json:"client_id" db:"client_id" gorm:"type:uuid;not null;index" validate:"required"`
	Scopes    string     `json:"scopes" db:"scopes" gorm:"type:text;not null" validate:"required"`
	GrantedAt time.Time  `json:"granted_at" db:"granted_at" gorm:"autoCreateTime"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	
	// Relationships
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Client Client `json:"client,omitempty" gorm:"foreignKey:ClientID"`
}

// AccountUser represents users who can manage accounts/tenants
type AccountUser struct {
	ID           uuid.UUID `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AccountID    uuid.UUID `json:"account_id" db:"account_id" gorm:"type:uuid;not null;index" validate:"required"`
	Email        string    `json:"email" db:"email" gorm:"type:varchar(255);not null" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash" gorm:"type:varchar(255);not null" validate:"required"`
	Role         string    `json:"role" db:"role" gorm:"type:varchar(50);not null" validate:"required,oneof=owner admin member"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Account   Account     `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	AuditLogs []AuditLog  `json:"audit_logs,omitempty" gorm:"foreignKey:AccountUserID;constraint:OnDelete:SET NULL"`
}

// AuditLog represents security and compliance audit trail
type AuditLog struct {
	ID            uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty" db:"tenant_id" gorm:"type:uuid;index"`
	UserID        *uuid.UUID `json:"user_id,omitempty" db:"user_id" gorm:"type:uuid;index"`
	AccountUserID *uuid.UUID `json:"account_user_id,omitempty" db:"account_user_id" gorm:"type:uuid;index"`
	Action        string     `json:"action" db:"action" gorm:"type:varchar(100);not null" validate:"required"`
	ResourceType  string     `json:"resource_type,omitempty" db:"resource_type" gorm:"type:varchar(50)"`
	ResourceID    *uuid.UUID `json:"resource_id,omitempty" db:"resource_id" gorm:"type:uuid"`
	IPAddress     string     `json:"ip_address,omitempty" db:"ip_address" gorm:"type:varchar(45)"`
	UserAgent     string     `json:"user_agent,omitempty" db:"user_agent" gorm:"type:text"`
	Metadata      []byte     `json:"metadata,omitempty" db:"metadata" gorm:"type:jsonb"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime;index"`
	
	// Relationships
	Tenant      *Tenant      `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	User        *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AccountUser *AccountUser `json:"account_user,omitempty" gorm:"foreignKey:AccountUserID"`
}

//TableName Overrides
func (Account) TableName() string          { return "accounts" }
func (Tenant) TableName() string           { return "tenants" }
func (Client) TableName() string           { return "clients" }
func (User) TableName() string             { return "users" }
func (AuthorizationCode) TableName() string{ return "authorization_codes" }
func (AccessToken) TableName() string      { return "access_tokens" }
func (RefreshToken) TableName() string     { return "refresh_tokens" }
func (IDToken) TableName() string          { return "id_tokens" }
func (Session) TableName() string          { return "sessions" }
func (UserConsent) TableName() string      { return "user_consents" }
func (AccountUser) TableName() string      { return "account_users" }
func (AuditLog) TableName() string         { return "audit_logs" }


//Tenant Functions
func (Tenant) CreateSlug() string {
	//word-word-numbers
	//3-4 letters-5-6 letters-4numbers
	//
	firstWord := shortWords[rand.IntN(len(shortWords))]
	secondWord := mediumWords[rand.IntN(len(mediumWords))]

	numbers := rand.IntN(9000) + 1000

	return fmt.Sprintf("%s-%s-%d", firstWord, secondWord, numbers)
}