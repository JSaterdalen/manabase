# Manabase
Product Requirements Document

## Overview
Manabase is a web application for tracking Magic: The Gathering games. Users can create decks, record games, and analyze their performance across different decks and opponents.

## Tech Stack
- Backend: Go
- Frontend: HTMX + TailwindCSS
- Templating: templ
- Database: PostgreSQL
- Query Builder: sqlc

## Project Structure
```
.
├── cmd/manabase
│   └── main.go                 # Application entry point
├── internal
│   ├── database               # Database layer
│   ├── handler                # HTTP handlers
│   └── middleware             # Auth middleware
├── sql                        # Database definitions
├── web                        # Frontend assets
└── sqlc.yaml                  # SQLC configuration
```

## Core Features

### 1. Database Schema
The database will use the following schema:

#### Users
- ID (UUID)
- Username (string, unique)
- Email (string, unique)
- PasswordHash (string)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

#### Players
- ID (UUID)
- UserID (UUID, foreign key)
- Name (string)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

#### Decks
- ID (UUID)
- OwnerID (UUID, foreign key)
- Name (string)
- Commander (string, nullable)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

#### Games
- ID (UUID)
- Date (timestamp)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)

#### PlayerDeckGame
- GameID (UUID, foreign key)
- PlayerID (UUID, foreign key)
- DeckID (UUID, foreign key)
- Primary Key (GameID, PlayerID)
- IsWon (boolean)

### 2. Authentication System

#### Requirements
- Session-based authentication using secure cookies
- Password hashing using bcrypt
- CSRF protection for all forms
- Rate limiting for login attempts

#### User Flows
1. Sign Up
   - Required fields: username, email, password
   - Validate email format
   - Ensure password strength (min 8 chars, 1 number, 1 special char)
   - Prevent duplicate usernames/emails

2. Sign In
   - Login with username/email + password
   - Remember me option
   - Redirect to last page after login

3. Sign Out
   - Clear session
   - Redirect to home page
   - Prevent session fixation

### 3. Home Page

#### Game List
- Show most recent games first
- Display:
  - Date
  - Players and their decks
  - Winner
  - Quick stats (total players, duration if tracked)
- Filter options:
  - By date range
  - By player
  - By deck

#### New Game Form
- Date picker (defaults to today)
- Player selection
  - Dynamic dropdown of user's players
  - Minimum 2 players required
- Deck selection
  - Show only decks owned by selected players
  - Required for each player
- Winner selection
  - Can only select from players in game
- Save & Add Another option

### 4. Deck Management

#### Deck List
- Display:
  - Deck name
  - Colors
  - Commander (if applicable)
  - Win rate
  - Total games
- Sort options:
  - By name
  - By win rate
  - By number of games

#### Deck Operations
1. Create
   - Required: name, player owner
   - Optional: colors, commander
   - Validate color combinations

2. Edit
   - Allow all fields to be updated
   - Track modification date

3. Delete
   - Require confirmation
   - Only if no games recorded
   - Archive instead if games exist

### 5. Player Management

#### Player List
- Show:
  - Player name
  - Number of decks
  - Total games
  - Win rate
- Sort by name/stats

#### Player Operations
1. Create
   - Required: name
   - Optional: preferred colors

2. Edit
   - Allow name update
   - Track modification date

3. Delete
   - Require confirmation
   - Only if no games/decks
   - Archive instead if data exists

### 6. Analytics

#### Game Analytics
- Win rates over time
- Number of games by month
- Popular matchups

#### Deck Analytics
- Win rate vs specific decks
- Win rate vs color combinations
- Performance over time

#### Player Analytics
- Overall win rate
- Favorite decks
- Best/worst matchups

## Technical Requirements

### Performance
- Page load < 1s
- Database queries < 100ms
- Support 1000 concurrent users

### Security
- HTTPS only
- Secure session management
- Input sanitization
- SQL injection prevention
- XSS protection

### Accessibility
- WCAG 2.1 AA compliance
- Keyboard navigation
- Screen reader support
- Color contrast requirements

### Browser Support
- Latest 2 versions of:
  - Chrome
  - Firefox
  - Safari
  - Edge

## Future Considerations
- Mobile app version
- Tournament tracking
- Deck import from external sources
- Social features (friends, playgroups)
- Integration with external APIs

## Success Metrics
- User engagement (games logged per user)
- User retention (30-day active users)
- System performance (response times)
- Error rates (< 0.1% of requests)