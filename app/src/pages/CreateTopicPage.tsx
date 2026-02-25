import { Link } from 'react-router-dom'
import CreateTopicForm from '../components/topics/CreateTopicForm'

export default function CreateTopicPage() {
  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Breadcrumb */}
      <Link
        to="/"
        className="text-sm text-primary-600 hover:text-primary-700 inline-block"
      >
        ← Back to all topics
      </Link>

      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Create New Topic</h1>
        <p className="text-gray-600">
          Start a new discussion by creating a topic. Give it a clear and descriptive title.
        </p>
      </div>

      {/* Form */}
      <div className="bg-white rounded-lg shadow-sm p-6 border border-gray-200">
        <CreateTopicForm />
      </div>
    </div>
  )
}
