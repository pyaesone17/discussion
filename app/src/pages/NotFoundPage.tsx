import { Link } from 'react-router-dom'
import Button from '../components/ui/Button'

export default function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center">
      <div className="space-y-6">
        <div className="text-9xl font-bold text-primary-600">404</div>
        <h1 className="text-3xl font-bold text-gray-900">Page Not Found</h1>
        <p className="text-gray-600 max-w-md">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <Link to="/">
          <Button>Go Back Home</Button>
        </Link>
      </div>
    </div>
  )
}
