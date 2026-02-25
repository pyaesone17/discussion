# Discussion Forum - Frontend

A modern React application for the Discussion Forum API, built with Vite, TypeScript, React Query, and Tailwind CSS.

## Features

- **Topic Management**: Browse, search, and create discussion topics
- **Post System**: View and reply to topics with threaded posts
- **User Selection**: Simple user authentication via localStorage
- **Real-time Updates**: Optimistic UI updates with React Query
- **Responsive Design**: Desktop-first design with mobile support
- **Toast Notifications**: User-friendly success/error notifications

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **React Query (TanStack Query)** - Server state management
- **Axios** - HTTP client
- **Tailwind CSS** - Utility-first styling

## Prerequisites

- Node.js 18+ and npm/yarn/pnpm
- Discussion Forum API running on `http://localhost:8080`

## Getting Started

### 1. Install Dependencies

```bash
npm install
```

### 2. Environment Setup

Create a `.env` file in the project root (copy from `.env.example`):

```bash
cp .env.example .env
```

The default configuration connects to the API at `http://localhost:8080/api/v1`.

### 3. Start Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:5173`.

### 4. Build for Production

```bash
npm run build
```

Built files will be in the `dist/` directory.

## Project Structure

```
src/
├── api/              # API client and service functions
│   ├── client.ts     # Axios instance with interceptors
│   ├── topics.ts     # Topic API calls
│   ├── posts.ts      # Post API calls
│   └── users.ts      # User API calls
│
├── components/       # React components
│   ├── layout/       # Layout components (Navbar, Layout)
│   ├── modals/       # Modal components (UserSelectModal)
│   ├── topics/       # Topic-related components
│   ├── posts/        # Post-related components
│   └── ui/           # Reusable UI primitives
│
├── context/          # React Context providers
│   ├── UserContext.tsx    # Current user state
│   └── ToastContext.tsx   # Toast notifications
│
├── hooks/            # Custom React hooks
│   └── useDebounce.ts     # Debounce hook for search
│
├── pages/            # Page components
│   ├── HomePage.tsx           # Topics list + search
│   ├── TopicDetailPage.tsx   # Topic + posts
│   ├── CreateTopicPage.tsx   # Create new topic
│   └── NotFoundPage.tsx      # 404 page
│
├── types/            # TypeScript type definitions
│   ├── user.ts       # User types
│   ├── topic.ts      # Topic types
│   └── post.ts       # Post types
│
├── utils/            # Utility functions
│   ├── constants.ts  # App constants
│   └── date.ts       # Date formatting
│
├── App.tsx           # Root component with routing
├── main.tsx          # App entry point
└── index.css         # Global styles + Tailwind
```

## Architecture

### Pragmatic Hybrid Architecture

This application uses a **pragmatic architecture** that balances simplicity with scalability:

- **React Query** handles all server state (topics, posts, users)
- **Context API** manages client state (selected user, toast notifications)
- **Axios** provides HTTP client with error interceptors
- **TypeScript** types mirror Go backend models for type safety

### Key Patterns

1. **Server State (React Query)**
   - Automatic caching with configurable TTLs
   - Optimistic updates on mutations
   - Auto-refetch and invalidation strategies

2. **Client State (Context API)**
   - User selection persisted to localStorage
   - Toast notifications with auto-dismiss

3. **API Layer**
   - Centralized Axios instance with error handling
   - Service functions for each resource (topics, posts, users)
   - Automatic error extraction from API responses

## Usage

### First Visit

1. On first visit, you'll be prompted to select a user
2. This selection is saved to localStorage and persists across sessions
3. All topics and posts you create will be attributed to this user

### Browse Topics

- View all topics on the home page
- Click any topic to view details and posts

### Search Topics

- Use the search bar in the navigation
- Search uses MySQL full-text search with relevance ranking
- Results are cached for 3 minutes (matches API cache)

### Create Topic

- Click "New Topic" button in navigation
- Enter a title (5-200 characters)
- Topic is created and you're redirected to the topic detail page

### Reply to Topics

- Open any topic detail page
- Scroll to the bottom to see the reply form
- Enter your message (1-5000 characters)
- Post is added immediately with optimistic UI update

## Development

### React Query DevTools

The app includes React Query DevTools for debugging. Press the React Query icon in the bottom-left corner to open the dev tools panel.

### API Error Handling

All API errors are extracted from the backend's `{error: "message"}` format and displayed via toast notifications.

### Caching Strategy

Caching configuration matches the backend's Redis cache TTLs:

- **Topics list**: 2 minutes
- **Single topic**: 5 minutes
- **Posts for topic**: 3 minutes
- **Search results**: 3 minutes

Cache is automatically invalidated on create/update/delete operations.

## Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Lint code with ESLint

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | API base URL | `http://localhost:8080/api/v1` |

## Troubleshooting

### API Connection Issues

If you see "Network error" or connection refused:

1. Ensure the Go API is running on port 8080
2. Check that `.env` has the correct `VITE_API_BASE_URL`
3. Verify CORS is enabled on the API for `http://localhost:5173`

### User Selection Modal Not Appearing

The modal only shows if no user is selected. To reset:

```javascript
// In browser console:
localStorage.removeItem('discussion_forum_user')
// Then refresh the page
```

### Build Errors

If you encounter TypeScript errors during build:

```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Contributing

1. Follow the existing code structure and patterns
2. Use TypeScript for all new files
3. Match Go backend types in TypeScript definitions
4. Add error handling for all API calls
5. Keep components small and focused

## License

This project is part of the Discussion Forum application.
