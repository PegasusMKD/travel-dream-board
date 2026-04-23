import { Routes, Route, Navigate } from 'react-router-dom'
import Header from './components/Header'
import Home from './pages/Home'
import BoardDetail from './pages/BoardDetail'
import Login from './pages/Login'
import { useAuth } from './context/AuthContext'
import { useShare } from './context/ShareContext'

function ProtectedRoute({ children }) {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen bg-surface-50 flex items-center justify-center">
        <div className="w-8 h-8 border-3 border-accent-200 border-t-accent-500 rounded-full animate-spin" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return children
}

function GuestOrProtectedBoard() {
  const { user, loading } = useAuth()
  const { shareToken } = useShare()

  if (loading) {
    return (
      <div className="min-h-screen bg-surface-50 flex items-center justify-center">
        <div className="w-8 h-8 border-3 border-accent-200 border-t-accent-500 rounded-full animate-spin" />
      </div>
    )
  }

  if (!user && !shareToken) {
    return <Navigate to="/login" replace />
  }

  return <BoardDetail />
}

export default function App() {
  const { user, loading } = useAuth()
  const { shareToken } = useShare()

  // Show login page when not authenticated AND no share token in scope
  if (!loading && !user && !shareToken) {
    return (
      <div className="min-h-screen bg-surface-50 text-text-primary transition-colors duration-200">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/board/:id" element={<GuestOrProtectedBoard />} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface-50 text-text-primary transition-colors duration-200">
      {user && <Header />}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16">
        <Routes>
          <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
          <Route path="/board/:id" element={<GuestOrProtectedBoard />} />
          <Route path="/login" element={<Navigate to="/" replace />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  )
}
