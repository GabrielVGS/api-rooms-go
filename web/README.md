# Room Booking Frontend

A React frontend application for the room booking system built with:

- **React 19** with TypeScript
- **Vite** for build tooling
- **shadcn/ui** component library
- **Tailwind CSS** for styling
- **React Router** for navigation
- **React Hook Form** with Zod validation
- **Axios** for API calls

## Quick Start

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Build for production:
   ```bash
   npm run build
   ```

## Features

- **Authentication**: Login and register with JWT tokens
- **Room Management**: Create, edit, and delete meeting rooms
- **Reservations**: Book rooms with date/time selection
- **Dashboard**: Overview of rooms and upcoming reservations
- **Responsive Design**: Works on desktop and mobile devices

## API Integration

The frontend expects the Go API to be running on `http://localhost:8080` with the following endpoints:

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET|POST|PUT|DELETE /api/rooms` - Room management
- `GET|POST|PUT|DELETE /api/reservations` - Reservation management

## Development

The application uses modern React patterns:

- **Context API** for authentication state
- **Custom hooks** for data fetching
- **TypeScript** for type safety
- **Form validation** with Zod schemas
- **Component composition** with shadcn/ui

## Theme

The application uses a custom shadcn/ui theme with purple/blue color scheme as specified in the project requirements.
