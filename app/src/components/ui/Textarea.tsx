import { TextareaHTMLAttributes } from 'react'

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  error?: string
}

export default function Textarea({ error, className = '', ...props }: TextareaProps) {
  const baseClasses = 'w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 transition-colors duration-200 resize-vertical'
  const errorClasses = error
    ? 'border-red-500 focus:ring-red-500'
    : 'border-gray-300 focus:ring-primary-500 focus:border-primary-500'

  return (
    <div className="w-full">
      <textarea
        className={`${baseClasses} ${errorClasses} ${className}`}
        {...props}
      />
      {error && (
        <p className="mt-1 text-sm text-red-600">{error}</p>
      )}
    </div>
  )
}
