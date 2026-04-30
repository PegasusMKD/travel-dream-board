import { useCallback, useEffect, useState } from 'react'
import {
  X,
  ExternalLink,
  Star,
  Bookmark,
  BookmarkX,
  Bed,
  Plane,
  MapPinned,
  Pencil,
  Save,
  XCircle,
  ImagePlus,
  Trash2,
  Send,
  Loader2,
  Check,
  ArrowRight,
  CornerUpLeft,
  Clock,
} from 'lucide-react'
import StatusBadge from './StatusBadge'
import { useLang } from '../context/LanguageContext'
import { useAuth } from '../context/AuthContext'
import { api } from '../services/api'
import {
  sectionToItemApi,
  toBackendItemPayload,
  backendCommentTarget,
  mapItem,
} from '../services/mappers'
import { itemsToUnselectFor, finalizeLimit, nextStatus } from '../services/finalizeRules'

const sectionIcons = {
  accommodation: Bed,
  transport: Plane,
  activities: MapPinned,
}

const ALL_STATUSES = ['considering', 'finalist', 'rejected', 'booked', 'completed']

export default function ItemDetailSidebar({
  itemId,
  sectionType,
  onClose,
  onRefresh,
  onVote,
  onClearVote,
  isGuest = false,
  voterUuid: voterUuidProp,
  guestName,
  board,
}) {
  const { t } = useLang()
  const { user } = useAuth()
  const voterUuid = voterUuidProp || user?.uuid || null

  const [item, setItem] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [editing, setEditing] = useState(false)
  const [saving, setSaving] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [commentText, setCommentText] = useState('')
  const [postingComment, setPostingComment] = useState(false)
  const [editingCommentId, setEditingCommentId] = useState(null)
  const [commentDraft, setCommentDraft] = useState('')
  const [commentBusyId, setCommentBusyId] = useState(null)

  const [draft, setDraft] = useState(() => buildDraft(null))

  const reloadItem = useCallback(async () => {
    const itemApi = sectionToItemApi(api, sectionType)
    const data = await itemApi.get(itemId)
    const mapped = mapItem(data)
    setItem(mapped)
    return mapped
  }, [itemId, sectionType])

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)
    setItem(null)
    setEditing(false)
    setCommentText('')
    setEditingCommentId(null)

    reloadItem()
      .then((mapped) => {
        if (cancelled) return
        setDraft(buildDraft(mapped))
      })
      .catch((err) => {
        if (!cancelled) setError(err.message)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => { cancelled = true }
  }, [reloadItem])

  const startEditing = () => {
    setDraft(buildDraft(item))
    setError(null)
    setEditing(true)
  }

  const cancelEditing = () => {
    setEditing(false)
    setError(null)
  }

  const handleSave = async () => {
    setSaving(true)
    setError(null)
    try {
      const itemApi = sectionToItemApi(api, sectionType)

      // If marking this item as final, enforce per-section limit by
      // unmarking conflicting items first.
      if (draft.isFinal && !item.isFinal && board?.sections?.[sectionType]) {
        const conflicts = itemsToUnselectFor(sectionType, board.sections[sectionType], item.id)
        for (const conflict of conflicts) {
          const conflictPayload = toBackendItemPayload({
            ...conflict,
            isFinal: false,
          })
          await itemApi.update(conflict.id, conflictPayload)
        }
      }

      const payload = toBackendItemPayload({
        id: item.id,
        boardUuid: item.boardUuid,
        url: draft.url,
        title: draft.title,
        image: draft.image,
        note: draft.note,
        status: draft.status,
        isFinal: draft.isFinal,
        bookingRef: draft.bookingRef,
        outboundDepartingLocation: draft.outboundDepartingLocation,
        outboundArrivingLocation: draft.outboundArrivingLocation,
        outboundDepartingAt: fromDatetimeLocalValue(draft.outboundDepartingAt),
        outboundArrivingAt: fromDatetimeLocalValue(draft.outboundArrivingAt),
        inboundDepartingLocation: draft.inboundDepartingLocation,
        inboundArrivingLocation: draft.inboundArrivingLocation,
        inboundDepartingAt: fromDatetimeLocalValue(draft.inboundDepartingAt),
        inboundArrivingAt: fromDatetimeLocalValue(draft.inboundArrivingAt),
        startAt: fromDatetimeLocalValue(draft.startAt),
        endAt: fromDatetimeLocalValue(draft.endAt),
        outboundDurationMinutes: parseDurationMinutes(draft.outboundDurationMinutes),
        inboundDurationMinutes: parseDurationMinutes(draft.inboundDurationMinutes),
      })
      await itemApi.update(item.id, payload)
      await reloadItem()
      await onRefresh?.()
      setEditing(false)
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  // Owner-only quick action: advance status by one step (considering → finalist → booked).
  const handleAdvanceStatus = async () => {
    if (!item || isGuest) return
    const next = nextStatus(item.status)
    if (!next) return
    if (next === 'booked' && !item.bookingRef) {
      // Booking requires a reference — open edit mode so the user can fill it in.
      setDraft({ ...buildDraft(item), status: 'booked' })
      setEditing(true)
      return
    }
    setSaving(true)
    setError(null)
    try {
      const itemApi = sectionToItemApi(api, sectionType)
      const payload = toBackendItemPayload({ ...item, status: next })
      await itemApi.update(item.id, payload)
      await reloadItem()
      await onRefresh?.()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  const handleAddComment = async (e) => {
    e.preventDefault()
    if (!commentText.trim() || !voterUuid || !item) return
    setPostingComment(true)
    setError(null)
    try {
      await api.comments.create({
        user_uuid: voterUuid,
        content: commentText.trim(),
        commented_on: backendCommentTarget(sectionType),
        commented_on_uuid: item.id,
      })
      setCommentText('')
      await reloadItem()
      await onRefresh?.()
    } catch (err) {
      setError(err.message)
    } finally {
      setPostingComment(false)
    }
  }

  const handleDeleteItem = async () => {
    if (!item) return
    if (!window.confirm(t.confirmDeleteItem)) return
    setDeleting(true)
    setError(null)
    try {
      const itemApi = sectionToItemApi(api, sectionType)
      await itemApi.delete(item.id)
      await onRefresh?.()
      onClose()
    } catch (err) {
      setError(err.message)
      setDeleting(false)
    }
  }

  const startEditingComment = (c) => {
    setEditingCommentId(c.id)
    setCommentDraft(c.text)
  }

  const cancelEditingComment = () => {
    setEditingCommentId(null)
    setCommentDraft('')
  }

  const handleSaveComment = async (c) => {
    if (!commentDraft.trim()) return
    setCommentBusyId(c.id)
    setError(null)
    try {
      await api.comments.update(c.id, commentDraft.trim())
      setEditingCommentId(null)
      setCommentDraft('')
      await reloadItem()
      await onRefresh?.()
    } catch (err) {
      setError(err.message)
    } finally {
      setCommentBusyId(null)
    }
  }

  const handleDeleteComment = async (c) => {
    if (!window.confirm(t.confirmDeleteComment)) return
    setCommentBusyId(c.id)
    setError(null)
    try {
      await api.comments.delete(c.id)
      await reloadItem()
      await onRefresh?.()
    } catch (err) {
      setError(err.message)
    } finally {
      setCommentBusyId(null)
    }
  }

  const updateDraft = (field, value) => setDraft((prev) => ({ ...prev, [field]: value }))

  const FallbackIcon = sectionIcons[sectionType] || MapPinned

  return (
    <>
      <div
        className="fixed inset-0 bg-black/30 backdrop-blur-sm z-40"
        onClick={onClose}
      />

      <div className="fixed top-0 right-0 h-full w-full max-w-md bg-surface-0 z-50 shadow-2xl flex flex-col animate-slide-in">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-surface-200">
          <h2 className="text-base font-bold text-text-primary">{t.details}</h2>
          <div className="flex items-center gap-1">
            {!loading && item && !editing && !isGuest && (
              <>
                <button
                  onClick={startEditing}
                  disabled={deleting}
                  className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
                >
                  <Pencil className="w-3.5 h-3.5" />
                  {t.edit}
                </button>
                <button
                  onClick={handleDeleteItem}
                  disabled={deleting}
                  title={t.deleteItem}
                  className="w-8 h-8 flex items-center justify-center rounded-lg text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors cursor-pointer disabled:opacity-50"
                >
                  {deleting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
                </button>
              </>
            )}
            {!loading && item && editing && (
              <>
                <button
                  onClick={cancelEditing}
                  disabled={saving}
                  className="flex items-center gap-1.5 text-xs font-semibold text-text-tertiary hover:text-text-secondary bg-surface-100 hover:bg-surface-200 px-3 py-1.5 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
                >
                  <XCircle className="w-3.5 h-3.5" />
                  {t.cancel}
                </button>
                <button
                  onClick={handleSave}
                  disabled={saving}
                  className="flex items-center gap-1.5 text-xs font-semibold text-white bg-accent-500 hover:bg-accent-600 px-3 py-1.5 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
                >
                  {saving ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Save className="w-3.5 h-3.5" />}
                  {t.save}
                </button>
              </>
            )}
            <button
              onClick={onClose}
              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer ml-1"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {loading ? (
          <div className="flex-1 flex items-center justify-center">
            <Loader2 className="w-6 h-6 text-accent-500 animate-spin" />
          </div>
        ) : !item ? (
          <div className="flex-1 flex items-center justify-center px-6 text-center">
            <p className="text-sm text-text-muted">{error || t.boardNotFound}</p>
          </div>
        ) : (
          <ItemSidebarBody
            item={item}
            sectionType={sectionType}
            editing={editing}
            draft={draft}
            updateDraft={updateDraft}
            error={error}
            user={user}
            voterUuid={voterUuid}
            isGuest={isGuest}
            t={t}
            FallbackIcon={FallbackIcon}
            commentText={commentText}
            setCommentText={setCommentText}
            postingComment={postingComment}
            handleAddComment={handleAddComment}
            editingCommentId={editingCommentId}
            commentDraft={commentDraft}
            setCommentDraft={setCommentDraft}
            commentBusyId={commentBusyId}
            startEditingComment={startEditingComment}
            cancelEditingComment={cancelEditingComment}
            handleSaveComment={handleSaveComment}
            handleDeleteComment={handleDeleteComment}
            onVote={onVote}
            onClearVote={onClearVote}
            reloadItem={reloadItem}
            handleAdvanceStatus={handleAdvanceStatus}
            saving={saving}
          />
        )}
      </div>

      <style>{`
        @keyframes slideIn {
          from { transform: translateX(100%); }
          to { transform: translateX(0); }
        }
        .animate-slide-in {
          animation: slideIn 0.25s ease-out;
        }
      `}</style>
    </>
  )
}

function ItemSidebarBody({
  item, sectionType, editing, draft, updateDraft, error, user, voterUuid, isGuest, t, FallbackIcon,
  commentText, setCommentText, postingComment, handleAddComment,
  editingCommentId, commentDraft, setCommentDraft, commentBusyId,
  startEditingComment, cancelEditingComment, handleSaveComment, handleDeleteComment,
  onVote, onClearVote, reloadItem, handleAdvanceStatus, saving,
}) {
  const displayImage = editing ? draft.image : item.image
  const displayTitle = editing ? draft.title : item.title
  const myVote = voterUuid ? item.votes.find((v) => v.userUuid === voterUuid) : null
  const sortedRaters = [...item.votes].sort((a, b) => b.rank - a.rank)

  const handleRate = async (rank) => {
    if (!voterUuid || !onVote) return
    try {
      await onVote(item, rank)
      await reloadItem()
    } catch {}
  }

  const handleClear = async () => {
    if (!voterUuid || !onClearVote) return
    try {
      await onClearVote(item)
      await reloadItem()
    } catch {}
  }

  return (
    <div className="flex-1 overflow-y-auto">
      {/* Image area */}
      <div className="relative group">
        {displayImage ? (
          <>
            <img src={displayImage} alt={displayTitle} className="w-full h-52 object-cover" />
            {editing && (
              <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-3">
                <button
                  onClick={() => {
                    const url = prompt(t.editImagePlaceholder)
                    if (url) updateDraft('image', url)
                  }}
                  className="flex items-center gap-1.5 bg-white text-gray-800 text-xs font-semibold px-3 py-2 rounded-lg cursor-pointer"
                >
                  <ImagePlus className="w-3.5 h-3.5" />
                  {t.changeImage}
                </button>
                <button
                  onClick={() => updateDraft('image', '')}
                  className="flex items-center gap-1.5 bg-red-500 text-white text-xs font-semibold px-3 py-2 rounded-lg cursor-pointer"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  {t.removeImage}
                </button>
              </div>
            )}
          </>
        ) : (
          <div className="w-full h-36 bg-surface-100 flex flex-col items-center justify-center gap-2">
            <FallbackIcon className="w-12 h-12 text-surface-300" />
            {editing && (
              <button
                onClick={() => {
                  const url = prompt(t.editImagePlaceholder)
                  if (url) updateDraft('image', url)
                }}
                className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
              >
                <ImagePlus className="w-3.5 h-3.5" />
                {t.changeImage}
              </button>
            )}
          </div>
        )}
      </div>

      <div className="p-6 space-y-5">
        {error && (
          <div className="text-xs text-red-500 bg-red-50 dark:bg-red-950/30 rounded-lg p-2.5">
            {error}
          </div>
        )}

        {/* Final badge + Title */}
        <div>
          {editing ? (
            <>
              <button
                onClick={() => updateDraft('isFinal', !draft.isFinal)}
                className={`inline-flex items-center gap-1.5 text-[11px] font-bold px-2.5 py-1 rounded-lg mb-3 cursor-pointer transition-colors ${
                  draft.isFinal
                    ? 'bg-accent-500 text-white'
                    : 'bg-surface-100 text-text-tertiary hover:bg-surface-200'
                }`}
              >
                {draft.isFinal ? (
                  <Bookmark className="w-3 h-3" fill="currentColor" />
                ) : (
                  <BookmarkX className="w-3 h-3" />
                )}
                {draft.isFinal ? t.unmarkAsFinal : t.markAsFinal}
              </button>

              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                {t.editTitle}
              </label>
              <input
                type="text"
                value={draft.title}
                onChange={(e) => updateDraft('title', e.target.value)}
                className="w-full text-lg font-bold text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 transition-colors"
              />
            </>
          ) : (
            <>
              {item.isFinal && (
                <span className="inline-flex items-center gap-1 bg-accent-500 text-white text-[11px] font-bold px-2.5 py-1 rounded-lg mb-2">
                  <Bookmark className="w-3 h-3" fill="currentColor" />
                  {t.selected}
                </span>
              )}
              <h3 className="text-lg font-bold text-text-primary leading-snug">
                {item.title}
              </h3>
            </>
          )}
        </div>

        {/* Status */}
        <div>
          <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
            {t.status}
          </label>
          {editing ? (
            <div className="flex flex-wrap gap-2">
              {ALL_STATUSES.map((s) => (
                <button
                  key={s}
                  onClick={() => updateDraft('status', s)}
                  className={`cursor-pointer rounded-full transition-all ${
                    draft.status === s
                      ? 'ring-2 ring-accent-400 ring-offset-2 ring-offset-surface-0'
                      : 'opacity-60 hover:opacity-100'
                  }`}
                >
                  <StatusBadge status={s} />
                </button>
              ))}
            </div>
          ) : (
            <div className="flex items-center justify-between gap-2">
              <StatusBadge status={item.status} />
              {!isGuest && nextStatus(item.status) && (
                <button
                  onClick={handleAdvanceStatus}
                  disabled={saving}
                  className="text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-2.5 py-1 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
                >
                  {nextStatus(item.status) === 'booked' && !item.bookingRef
                    ? t.markBooked
                    : `→ ${t[nextStatus(item.status)]}`}
                </button>
              )}
            </div>
          )}
        </div>

        {/* Transport legs — outbound + optional inbound */}
        {sectionType === 'transport' && (
          editing ? (
            <div className="space-y-4">
              <TransportLegEditor
                title={t.transportOutbound}
                icon={Plane}
                fromValue={draft.outboundDepartingLocation}
                toValue={draft.outboundArrivingLocation}
                departValue={draft.outboundDepartingAt}
                arriveValue={draft.outboundArrivingAt}
                durationValue={draft.outboundDurationMinutes}
                onFrom={(v) => updateDraft('outboundDepartingLocation', v)}
                onTo={(v) => updateDraft('outboundArrivingLocation', v)}
                onDepart={(v) => updateDraft('outboundDepartingAt', v)}
                onArrive={(v) => updateDraft('outboundArrivingAt', v)}
                onDuration={(v) => updateDraft('outboundDurationMinutes', v)}
                t={t}
              />
              <TransportLegEditor
                title={t.transportInbound}
                icon={CornerUpLeft}
                hint={t.transportInboundHint}
                fromValue={draft.inboundDepartingLocation}
                toValue={draft.inboundArrivingLocation}
                departValue={draft.inboundDepartingAt}
                arriveValue={draft.inboundArrivingAt}
                durationValue={draft.inboundDurationMinutes}
                onFrom={(v) => updateDraft('inboundDepartingLocation', v)}
                onTo={(v) => updateDraft('inboundArrivingLocation', v)}
                onDepart={(v) => updateDraft('inboundDepartingAt', v)}
                onArrive={(v) => updateDraft('inboundArrivingAt', v)}
                onDuration={(v) => updateDraft('inboundDurationMinutes', v)}
                t={t}
              />
            </div>
          ) : (
            hasAnyLegData(item) && (
              <div className="space-y-2">
                <LegSummary
                  icon={Plane}
                  label={t.transportOutbound}
                  from={item.outboundDepartingLocation}
                  to={item.outboundArrivingLocation}
                  departAt={item.outboundDepartingAt}
                  arriveAt={item.outboundArrivingAt}
                  durationMinutes={item.outboundDurationMinutes}
                  t={t}
                />
                <LegSummary
                  icon={CornerUpLeft}
                  label={t.transportInbound}
                  from={item.inboundDepartingLocation}
                  to={item.inboundArrivingLocation}
                  departAt={item.inboundDepartingAt}
                  arriveAt={item.inboundArrivingAt}
                  durationMinutes={item.inboundDurationMinutes}
                  t={t}
                />
              </div>
            )
          )
        )}

        {/* Activity start + end times */}
        {sectionType === 'activities' && (
          editing ? (
            <ActivityTimeEditor
              startValue={draft.startAt}
              endValue={draft.endAt}
              onStart={(v) => updateDraft('startAt', v)}
              onEnd={(v) => updateDraft('endAt', v)}
              t={t}
            />
          ) : (
            (item.startAt || item.endAt) && (
              <ActivityTimeSummary startAt={item.startAt} endAt={item.endAt} t={t} />
            )
          )
        )}

        {/* Booking ref — only relevant when status is booked */}
        {editing && draft.status === 'booked' ? (
          <div>
            <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
              {t.editBookingRef}
            </label>
            <input
              type="text"
              value={draft.bookingRef}
              onChange={(e) => updateDraft('bookingRef', e.target.value)}
              placeholder={t.editBookingRefPlaceholder}
              className="w-full text-sm font-mono text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
            />
            {!draft.bookingRef && (
              <p className="text-[11px] text-text-muted mt-1">{t.bookingRefRequiredHint}</p>
            )}
          </div>
        ) : !editing && item.status === 'booked' && item.bookingRef ? (
          <div className="bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-400 rounded-xl p-3 text-sm font-mono">
            {t.bookingRef}: <strong>{item.bookingRef}</strong>
          </div>
        ) : null}

        {/* Notes */}
        <div>
          <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
            {t.notes}
          </label>
          {editing ? (
            <textarea
              value={draft.note}
              onChange={(e) => updateDraft('note', e.target.value)}
              placeholder={t.editNotesPlaceholder}
              rows={3}
              className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted resize-y leading-relaxed transition-colors"
            />
          ) : item.note ? (
            <p className="text-sm text-text-secondary bg-surface-50 rounded-xl p-3 leading-relaxed">
              {item.note}
            </p>
          ) : (
            <p className="text-sm text-text-muted italic">{t.noNotes}</p>
          )}
        </div>

        {/* URL */}
        {editing ? (
          <div>
            <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
              {t.editUrl}
            </label>
            <input
              type="url"
              value={draft.url}
              onChange={(e) => updateDraft('url', e.target.value)}
              className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
            />
          </div>
        ) : (
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-sm font-semibold text-accent-500 hover:text-accent-600 transition-colors no-underline"
          >
            <ExternalLink className="w-4 h-4" />
            {t.openLink}
          </a>
        )}

        {editing && (
          <div>
            <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
              {t.editImageUrl}
            </label>
            <input
              type="url"
              value={draft.image}
              onChange={(e) => updateDraft('image', e.target.value)}
              placeholder={t.editImagePlaceholder}
              className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
            />
          </div>
        )}

        <div className="border-t border-surface-200" />

        {/* Comments */}
        <div>
          <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-2 block">
            {t.comments} ({item.comments.length})
          </label>
          {item.comments.length > 0 ? (
            <div className="space-y-3">
              {item.comments.map((c) => {
                const isOwn = voterUuid && c.userUuid === voterUuid
                const isEditingThis = editingCommentId === c.id
                const isBusy = commentBusyId === c.id
                return (
                  <div key={c.id} className="bg-surface-50 rounded-xl p-3 group">
                    <div className="flex items-center gap-2 mb-1.5">
                      <span className="w-6 h-6 rounded-full bg-accent-100 flex items-center justify-center text-[10px] font-bold text-accent-600">
                        {c.displayName[0]}
                      </span>
                      <span className="text-xs font-semibold text-text-primary">
                        {c.displayName}
                      </span>
                      {isOwn && !isEditingThis && (
                        <div className="ml-auto flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button
                            onClick={() => startEditingComment(c)}
                            disabled={isBusy}
                            title={t.edit}
                            className="w-6 h-6 flex items-center justify-center rounded text-text-tertiary hover:text-accent-500 hover:bg-accent-50 transition-colors cursor-pointer disabled:opacity-50"
                          >
                            <Pencil className="w-3 h-3" />
                          </button>
                          <button
                            onClick={() => handleDeleteComment(c)}
                            disabled={isBusy}
                            title={t.delete}
                            className="w-6 h-6 flex items-center justify-center rounded text-text-tertiary hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors cursor-pointer disabled:opacity-50"
                          >
                            {isBusy ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
                          </button>
                        </div>
                      )}
                    </div>
                    {isEditingThis ? (
                      <div className="pl-8 space-y-2">
                        <textarea
                          value={commentDraft}
                          onChange={(e) => setCommentDraft(e.target.value)}
                          rows={2}
                          autoFocus
                          className="w-full text-sm text-text-primary bg-surface-0 border border-surface-200 rounded-lg px-3 py-2 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 resize-y transition-colors"
                        />
                        <div className="flex items-center gap-2 justify-end">
                          <button
                            onClick={cancelEditingComment}
                            disabled={isBusy}
                            className="flex items-center gap-1 text-xs font-semibold text-text-tertiary hover:text-text-secondary px-2 py-1 rounded transition-colors cursor-pointer disabled:opacity-50"
                          >
                            <XCircle className="w-3 h-3" />
                            {t.cancel}
                          </button>
                          <button
                            onClick={() => handleSaveComment(c)}
                            disabled={isBusy || !commentDraft.trim()}
                            className="flex items-center gap-1 text-xs font-semibold text-white bg-accent-500 hover:bg-accent-600 px-2.5 py-1 rounded transition-colors cursor-pointer disabled:opacity-50"
                          >
                            {isBusy ? <Loader2 className="w-3 h-3 animate-spin" /> : <Check className="w-3 h-3" />}
                            {t.save}
                          </button>
                        </div>
                      </div>
                    ) : (
                      <p className="text-sm text-text-secondary leading-relaxed pl-8">
                        {c.text}
                      </p>
                    )}
                  </div>
                )
              })}
            </div>
          ) : (
            <p className="text-sm text-text-muted italic">{t.noComments}</p>
          )}

          <form onSubmit={handleAddComment} className="flex items-center gap-2 mt-3">
            <input
              type="text"
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
              placeholder={t.addComment}
              disabled={postingComment}
              className="flex-1 text-sm bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted text-text-primary transition-colors disabled:opacity-50"
            />
            <button
              type="submit"
              disabled={!commentText.trim() || postingComment}
              className="w-10 h-10 flex items-center justify-center bg-accent-500 hover:bg-accent-600 text-white rounded-xl transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {postingComment ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
            </button>
          </form>
        </div>

        <div className="border-t border-surface-200" />

        {/* Ratings */}
        <div>
          <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-2 block">
            {t.ratings} ({item.ratingCount ?? 0})
          </label>

          {voterUuid ? (
            <div className="bg-surface-50 rounded-xl p-3 mb-3">
              <div className="flex items-center justify-between">
                <div className="flex flex-col gap-1">
                  <span className="text-[11px] font-semibold text-text-tertiary uppercase tracking-wider">
                    {t.yourRating}
                  </span>
                  <StarPicker value={myVote?.rank ?? 0} onChange={handleRate} />
                </div>
                {myVote && (
                  <button
                    onClick={handleClear}
                    className="text-xs font-semibold text-text-tertiary hover:text-red-500 transition-colors cursor-pointer"
                  >
                    {t.clearRating}
                  </button>
                )}
              </div>
            </div>
          ) : null}

          {sortedRaters.length > 0 ? (
            <div className="space-y-1.5">
              {sortedRaters.map((v) => (
                <div key={v.id} className="flex items-center gap-2">
                  <span className="w-6 h-6 rounded-full bg-accent-100 flex items-center justify-center text-[10px] font-bold text-accent-600 shrink-0">
                    {v.displayName[0]}
                  </span>
                  <span className="text-xs font-semibold text-text-primary flex-1 truncate">
                    {v.displayName}
                  </span>
                  <StarDisplay value={v.rank} />
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-text-muted italic">{t.noRatingsYet}</p>
          )}
        </div>
      </div>
    </div>
  )
}

function buildDraft(item) {
  return {
    title: item?.title || '',
    note: item?.note || '',
    url: item?.url || '',
    image: item?.image || '',
    status: item?.status || 'considering',
    isFinal: !!item?.isFinal,
    bookingRef: item?.bookingRef || '',
    outboundDepartingLocation: item?.outboundDepartingLocation || '',
    outboundArrivingLocation: item?.outboundArrivingLocation || '',
    outboundDepartingAt: toDatetimeLocalValue(item?.outboundDepartingAt),
    outboundArrivingAt: toDatetimeLocalValue(item?.outboundArrivingAt),
    inboundDepartingLocation: item?.inboundDepartingLocation || '',
    inboundArrivingLocation: item?.inboundArrivingLocation || '',
    inboundDepartingAt: toDatetimeLocalValue(item?.inboundDepartingAt),
    inboundArrivingAt: toDatetimeLocalValue(item?.inboundArrivingAt),
    startAt: toDatetimeLocalValue(item?.startAt),
    endAt: toDatetimeLocalValue(item?.endAt),
    outboundDurationMinutes: item?.outboundDurationMinutes != null ? String(item.outboundDurationMinutes) : '',
    inboundDurationMinutes: item?.inboundDurationMinutes != null ? String(item.inboundDurationMinutes) : '',
  }
}

// Times are stored as "floating" wall-clock values: the ISO string is marked UTC
// but the Y/M/D H:M components are the literal clock reading at the location
// (e.g., 14:30 at the airport stays 14:30 regardless of browser TZ).
function toDatetimeLocalValue(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getUTCFullYear()}-${pad(d.getUTCMonth() + 1)}-${pad(d.getUTCDate())}T${pad(d.getUTCHours())}:${pad(d.getUTCMinutes())}`
}

function fromDatetimeLocalValue(s) {
  if (!s) return null
  // Don't use new Date().toISOString() — that shifts by browser TZ.
  // Treat the raw input components as wall-clock, mark as UTC.
  const m = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})(?::(\d{2}))?$/.exec(s)
  if (!m) return null
  const [, y, mo, d, h, mi, se] = m
  return `${y}-${mo}-${d}T${h}:${mi}:${se || '00'}.000Z`
}

function hasAnyLegData(item) {
  if (!item) return false
  return !!(
    item.outboundDepartingLocation ||
    item.outboundArrivingLocation ||
    item.outboundDepartingAt ||
    item.outboundArrivingAt ||
    item.inboundDepartingLocation ||
    item.inboundArrivingLocation ||
    item.inboundDepartingAt ||
    item.inboundArrivingAt ||
    item.outboundDurationMinutes != null ||
    item.inboundDurationMinutes != null
  )
}

function parseDurationMinutes(s) {
  if (s == null || s === '') return null
  const n = parseInt(s, 10)
  if (Number.isNaN(n) || n < 0) return null
  return n
}

function formatDurationMinutes(mins) {
  if (mins == null || mins < 0) return null
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h && m) return `${h}h ${m}m`
  if (h) return `${h}h`
  return `${m}m`
}

function TransportLegEditor({
  title, icon: Icon, hint, fromValue, toValue, departValue, arriveValue, durationValue,
  onFrom, onTo, onDepart, onArrive, onDuration, t,
}) {
  const inputClass = 'w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors'
  const labelClass = 'text-[11px] font-semibold text-text-tertiary uppercase tracking-wider mb-1 block'
  return (
    <div className="border border-surface-200 rounded-xl p-3 space-y-3">
      <div className="flex items-center gap-2">
        <Icon className="w-4 h-4 text-accent-500" />
        <span className="text-xs font-bold text-text-primary uppercase tracking-wider">{title}</span>
        {hint && <span className="text-[11px] text-text-muted ml-auto italic">{hint}</span>}
      </div>
      <div className="grid grid-cols-2 gap-2">
        <div>
          <label className={labelClass}>{t.transportFrom}</label>
          <input type="text" value={fromValue} onChange={(e) => onFrom(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label className={labelClass}>{t.transportTo}</label>
          <input type="text" value={toValue} onChange={(e) => onTo(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label className={labelClass}>{t.transportDeparture}</label>
          <input type="datetime-local" value={departValue} onChange={(e) => onDepart(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label className={labelClass}>{t.transportArrival}</label>
          <input type="datetime-local" value={arriveValue} onChange={(e) => onArrive(e.target.value)} className={inputClass} />
        </div>
        <div className="col-span-2">
          <label className={labelClass}>{t.transportDuration}</label>
          <input
            type="number"
            min="0"
            value={durationValue}
            onChange={(e) => onDuration(e.target.value)}
            placeholder={t.transportDurationPlaceholder}
            className={inputClass}
          />
        </div>
      </div>
    </div>
  )
}

function LegSummary({ icon: Icon, label, from, to, departAt, arriveAt, durationMinutes, t }) {
  if (!from && !to && !departAt && !arriveAt && durationMinutes == null) return null

  const formatWhen = (iso) => {
    if (!iso) return null
    const d = new Date(iso)
    if (isNaN(d.getTime())) return null
    return d.toLocaleString(undefined, {
      day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit', timeZone: 'UTC',
    })
  }

  const depart = formatWhen(departAt)
  const arrive = formatWhen(arriveAt)
  const duration = formatDurationMinutes(durationMinutes)

  return (
    <div className="flex items-start gap-2.5 text-sm bg-surface-50 rounded-xl px-3 py-2.5">
      <Icon className="w-4 h-4 text-accent-500 mt-0.5 shrink-0" />
      <div className="min-w-0 flex-1">
        <div className="text-[11px] font-semibold text-text-tertiary uppercase tracking-wider mb-0.5">{label}</div>
        {(from || to || duration) && (
          <div className="flex items-center gap-1.5 text-text-primary font-medium">
            <span className="truncate">{from || '—'}</span>
            <ArrowRight className="w-3.5 h-3.5 text-text-muted shrink-0" />
            <span className="truncate">{to || '—'}</span>
            {duration && (
              <span className="ml-auto text-xs font-normal text-text-tertiary shrink-0">
                {duration}
              </span>
            )}
          </div>
        )}
        {(depart || arrive) && (
          <div className="flex items-center gap-1.5 text-text-secondary text-xs mt-0.5">
            <span>{depart || '—'}</span>
            <ArrowRight className="w-3 h-3 text-text-muted shrink-0" />
            <span>{arrive || '—'}</span>
          </div>
        )}
      </div>
    </div>
  )
}

function ActivityTimeEditor({ startValue, endValue, onStart, onEnd, t }) {
  const inputClass = 'w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors'
  const labelClass = 'text-[11px] font-semibold text-text-tertiary uppercase tracking-wider mb-1 block'
  return (
    <div className="border border-surface-200 rounded-xl p-3 space-y-3">
      <div className="flex items-center gap-2">
        <Clock className="w-4 h-4 text-accent-500" />
        <span className="text-xs font-bold text-text-primary uppercase tracking-wider">{t.activityWhen}</span>
      </div>
      <div className="grid grid-cols-2 gap-2">
        <div>
          <label className={labelClass}>{t.activityStart}</label>
          <input type="datetime-local" value={startValue} onChange={(e) => onStart(e.target.value)} className={inputClass} />
        </div>
        <div>
          <label className={labelClass}>{t.activityEnd}</label>
          <input type="datetime-local" value={endValue} onChange={(e) => onEnd(e.target.value)} className={inputClass} />
        </div>
      </div>
    </div>
  )
}

function ActivityTimeSummary({ startAt, endAt, t }) {
  const formatWhen = (iso) => {
    if (!iso) return null
    const d = new Date(iso)
    if (isNaN(d.getTime())) return null
    return d.toLocaleString(undefined, {
      day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit', timeZone: 'UTC',
    })
  }
  const start = formatWhen(startAt)
  const end = formatWhen(endAt)
  return (
    <div className="flex items-start gap-2.5 text-sm bg-surface-50 rounded-xl px-3 py-2.5">
      <Clock className="w-4 h-4 text-accent-500 mt-0.5 shrink-0" />
      <div className="min-w-0 flex-1">
        <div className="text-[11px] font-semibold text-text-tertiary uppercase tracking-wider mb-0.5">{t.activityWhen}</div>
        <div className="flex items-center gap-1.5 text-text-primary font-medium">
          <span>{start || '—'}</span>
          <ArrowRight className="w-3.5 h-3.5 text-text-muted shrink-0" />
          <span>{end || '—'}</span>
        </div>
      </div>
    </div>
  )
}

function StarPicker({ value, onChange }) {
  const [hover, setHover] = useState(0)
  const active = hover || value
  return (
    <div className="flex items-center gap-0.5" onMouseLeave={() => setHover(0)}>
      {[1, 2, 3, 4, 5].map((n) => (
        <button
          key={n}
          type="button"
          onMouseEnter={() => setHover(n)}
          onClick={() => onChange(n)}
          className="p-0.5 cursor-pointer transition-transform hover:scale-110"
          title={`${n} / 5`}
        >
          <Star
            className={`w-5 h-5 ${n <= active ? 'text-amber-400' : 'text-surface-300'}`}
            fill={n <= active ? 'currentColor' : 'none'}
          />
        </button>
      ))}
    </div>
  )
}

function StarDisplay({ value }) {
  return (
    <div className="flex items-center gap-0.5 shrink-0">
      {[1, 2, 3, 4, 5].map((n) => (
        <Star
          key={n}
          className={`w-3.5 h-3.5 ${n <= value ? 'text-amber-400' : 'text-surface-300'}`}
          fill={n <= value ? 'currentColor' : 'none'}
        />
      ))}
    </div>
  )
}
