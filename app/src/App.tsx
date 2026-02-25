import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/layout/Layout'
import HomePage from './pages/HomePage'
import TopicDetailPage from './pages/TopicDetailPage'
import CreateTopicPage from './pages/CreateTopicPage'
import NotFoundPage from './pages/NotFoundPage'
import UserSelectModal from './components/modals/UserSelectModal'

function App() {
  return (
    <BrowserRouter>
      <UserSelectModal />
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="topics/new" element={<CreateTopicPage />} />
          <Route path="topics/:id" element={<TopicDetailPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
