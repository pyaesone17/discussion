import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { createReport } from '../../api/reports'
import { useUser } from '../../context/UserContext'
import { useToast } from '../../context/ToastContext'
import Modal from './Modal'
import Button from './Button'

const REPORT_REASONS = [
  { value: 'spam', label: 'Spam' },
  { value: 'harassment', label: 'Harassment' },
  { value: 'off_topic', label: 'Off Topic' },
  { value: 'misinformation', label: 'Misinformation' },
  { value: 'other', label: 'Other' },
]

interface ReportButtonProps {
  reportableType: 'topic' | 'post'
  reportableId: number
}

export default function ReportButton({
  reportableType,
  reportableId,
}: ReportButtonProps) {
  const { currentUser } = useUser()
  const { showToast } = useToast()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedReason, setSelectedReason] = useState<string>('')

  const reportMutation = useMutation({
    mutationFn: async (reason: string) => {
      if (!currentUser) throw new Error('No user selected')
      await createReport(reportableType, reportableId, currentUser.id, reason)
    },
    onSuccess: () => {
      showToast('Report submitted successfully', 'success')
      setIsModalOpen(false)
      setSelectedReason('')
    },
    onError: (error: Error) => {
      showToast(error.message, 'error')
    },
  })

  const handleOpenModal = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!currentUser) {
      showToast('Please select a user first', 'error')
      return
    }
    setIsModalOpen(true)
  }

  const handleSubmit = () => {
    if (!selectedReason) {
      showToast('Please select a reason', 'error')
      return
    }
    reportMutation.mutate(selectedReason)
  }

  return (
    <>
      <button
        onClick={handleOpenModal}
        className="p-1 rounded text-gray-400 hover:text-red-500 transition-colors"
        aria-label={`Report this ${reportableType}`}
        title={`Report this ${reportableType}`}
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2z"
          />
        </svg>
      </button>

      <Modal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false)
          setSelectedReason('')
        }}
        title={`Report ${reportableType}`}
      >
        <div className="space-y-4">
          <p className="text-sm text-gray-600">
            Why are you reporting this {reportableType}?
          </p>
          <div className="space-y-2">
            {REPORT_REASONS.map((reason) => (
              <label
                key={reason.value}
                className={`flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors ${
                  selectedReason === reason.value
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:bg-gray-50'
                }`}
              >
                <input
                  type="radio"
                  name="report-reason"
                  value={reason.value}
                  checked={selectedReason === reason.value}
                  onChange={(e) => setSelectedReason(e.target.value)}
                  className="text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm font-medium text-gray-700">
                  {reason.label}
                </span>
              </label>
            ))}
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                setIsModalOpen(false)
                setSelectedReason('')
              }}
            >
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleSubmit}
              isLoading={reportMutation.isPending}
              disabled={!selectedReason}
            >
              Submit Report
            </Button>
          </div>
        </div>
      </Modal>
    </>
  )
}
