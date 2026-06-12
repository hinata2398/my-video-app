"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import HlsPlayer from "../../components/HlsPlayer";
import BookmarkButton from "../../components/BookmarkButton";

type Video = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  video_url: string;
  status: string;
  view_count: number;
  created_at: string;
};

type Comment = {
  id: number;
  user_id: number;
  username: string;
  body: string;
  like_count: number;
  liked: boolean;
  dislike_count: number;
  disliked: boolean;
  created_at: string;
};

export default function VideoDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [liked, setLiked] = useState(false);
  const [dislikeCount, setDislikeCount] = useState(0);
  const [disliked, setDisliked] = useState(false);
  const [bookmarked, setBookmarked] = useState(false);
  const [comments, setComments] = useState<Comment[]>([]);
  const [commentBody, setCommentBody] = useState("");
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  const fetchVideo = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`)
      .then((res) => {
        if (!res.ok) {
          setNotFound(true);
          return null;
        }
        return res.json();
      })
      .then((data) => {
        if (data) setVideo(data);
      })
      .catch(() => setNotFound(true))
      .finally(() => setLoading(false));
  };
  const incrementViewCount = () => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/view`, {
      method: "POST",
    });
  };
  const fetchLikes = async () => {
    const [likeRes, dislikeRes] = await Promise.all([
      fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/likes`),
      fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/dislikes`),
    ]);
    const likeData = await likeRes.json();
    const dislikeData = await dislikeRes.json();
    setLikeCount(likeData.count);
    setDislikeCount(dislikeData.count);

    // ログイン中なら自分のいいね・よくないね状態を取得
    const token = localStorage.getItem("token");
    if (token) {
      const statusRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/like-status`,
        { headers: { Authorization: `Bearer ${token}` } },
      );
      if (statusRes.ok) {
        const statusData = await statusRes.json();
        setLiked(statusData.liked);
        setDisliked(statusData.disliked);
      }
    }
  };

  const fetchBookmarkStatus = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/bookmark-status`,
      { headers: { Authorization: `Bearer ${token}` } },
    );
    if (res.ok) {
      const data = await res.json();
      setBookmarked(data.bookmarked);
    }
  };

  const handleLike = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      // 未ログインなら認証ページへ
      router.push("/auth");
      return;
    }

    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/like`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const data = await res.json();

    setLiked(data.liked);
    setLikeCount(data.count);
    if (data.liked && disliked) {
      setDisliked(false);
      setDislikeCount((prev) => Math.max(0, prev - 1));
    }
  };

  const handleDislike = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/auth");
      return;
    }

    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/dislike`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const data = await res.json();
    setDisliked(data.disliked);
    setDislikeCount(data.count);
    if (data.disliked && liked) {
      setLiked(false);
      setLikeCount((prev) => Math.max(0, prev - 1));
    }
  };

  const fetchComments = async () => {
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/comments`,
    );
    const data = await res.json();
    setComments(data);
  };

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/auth");
      return;
    }
    if (!commentBody.trim()) return;

    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/comments`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ body: commentBody }),
      },
    );
    if (res.ok) {
      setCommentBody("");
      fetchComments();
    }
  };

  const handleCommentLike = async (commentId: number) => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/auth");
      return;
    }
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/comments/${commentId}/like`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    if (res.ok) {
      const data = await res.json();
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId
            ? {
                ...c,
                like_count: data.count,
                liked: data.liked,
                dislike_count: data.liked
                  ? Math.max(0, c.dislike_count - 1)
                  : c.dislike_count,
                disliked: data.liked ? false : c.disliked,
              }
            : c,
        ),
      );
    }
  };

  const handleCommentDislike = async (commentId: number) => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/auth");
      return;
    }
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/comments/${commentId}/dislike`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    if (res.ok) {
      const data = await res.json();
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId
            ? {
                ...c,
                dislike_count: data.count,
                disliked: data.disliked,
                like_count: data.disliked
                  ? Math.max(0, c.like_count - 1)
                  : c.like_count,
                liked: data.disliked ? false : c.liked,
              }
            : c,
        ),
      );
    }
  };

  // useEffect でページ読み込み時に呼ぶ
  useEffect(() => {
    fetchVideo();
    fetchLikes();
    fetchBookmarkStatus();
    fetchComments();
  }, [id]);

  // status が pending/processing の間はポーリング
  useEffect(() => {
    if (!video) return;
    if (video.status === "pending" || video.status === "processing") {
      intervalRef.current = setInterval(fetchVideo, 3000);
    } else {
      if (intervalRef.current) clearInterval(intervalRef.current);
    }
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [video?.status]);

  if (loading)
    return (
      <main style={{ padding: "2rem" }}>
        <p style={{ color: "#888" }}>読み込み中...</p>
      </main>
    );

  if (notFound || !video)
    return (
      <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem" }}>
        <p>動画が見つかりません。</p>
        <Link href="/" style={{ color: "#e00", textDecoration: "none" }}>
          ← 一覧へ戻る
        </Link>
      </main>
    );

  const isProcessing =
    video.status === "pending" || video.status === "processing";

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem" }}>
      <Link href="/" style={{ color: "#888", textDecoration: "none" }}>
        ← 一覧へ戻る
      </Link>

      {/* ステータスバナー */}
      {isProcessing && (
        <div
          style={{
            marginTop: "1rem",
            padding: "0.75rem 1rem",
            background: "#fff8e1",
            border: "1px solid #ffe082",
            borderRadius: 6,
            color: "#795548",
          }}
        >
          ⏳ 動画を変換中です。しばらくすると再生できるようになります...
        </div>
      )}
      {video.status === "error" && (
        <div
          style={{
            marginTop: "1rem",
            padding: "0.75rem 1rem",
            background: "#ffebee",
            border: "1px solid #ef9a9a",
            borderRadius: 6,
            color: "#c62828",
          }}
        >
          ⚠️ 動画の変換に失敗しました
        </div>
      )}

      <div
        style={{
          marginTop: "1rem",
          borderRadius: 8,
          overflow: "hidden",
          background: "#111",
        }}
      >
        {video.video_url && !isProcessing ? (
          <HlsPlayer
            src={video.video_url}
            poster={video.thumbnail_url || undefined}
            onView={incrementViewCount}
          />
        ) : video.thumbnail_url ? (
          <img
            src={video.thumbnail_url}
            alt={video.title}
            style={{
              width: "100%",
              maxHeight: 450,
              objectFit: "cover",
              display: "block",
            }}
          />
        ) : (
          <div
            style={{
              height: 400,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <span style={{ color: "#555", fontSize: "4rem" }}>▶</span>
          </div>
        )}
      </div>

      <div
        style={{
          marginTop: "1rem",
          display: "flex",
          alignItems: "flex-start",
          gap: "1rem",
        }}
      >
        <div>
          <h1 style={{ margin: "0 0 0.25rem" }}>{video.title}</h1>
          <p style={{ color: "#888", fontSize: "0.9rem", margin: 0 }}>
            {new Date(video.created_at).toLocaleDateString("ja-JP")}　
            {video.view_count}回視聴
          </p>
        </div>
        <div style={{ marginLeft: "auto", display: "flex", gap: 8 }}>
          <button
            onClick={handleLike}
            style={{
              background: "none",
              border: `1px solid ${liked ? "#e00" : "#ddd"}`,
              borderRadius: 20,
              cursor: "pointer",
              color: liked ? "#e00" : "#888",
              fontSize: "0.9rem",
              padding: "5px 14px",
              display: "flex",
              alignItems: "center",
              gap: 4,
            }}
          >
            👍 {likeCount}
          </button>
          <button
            onClick={handleDislike}
            style={{
              background: "none",
              border: `1px solid ${disliked ? "#555" : "#ddd"}`,
              borderRadius: 20,
              cursor: "pointer",
              color: disliked ? "#555" : "#888",
              fontSize: "0.9rem",
              padding: "5px 14px",
              display: "flex",
              alignItems: "center",
              gap: 4,
            }}
          >
            👎 {dislikeCount}
          </button>
          <BookmarkButton videoID={video.id} initialBookmarked={bookmarked} />
        </div>
      </div>

      {video.description && (
        <p
          style={{
            marginTop: "1rem",
            lineHeight: 1.7,
            borderTop: "1px solid #eee",
            paddingTop: "1rem",
          }}
        >
          {video.description}
        </p>
      )}

      {/* コメントセクション */}
      <div
        style={{
          marginTop: "2rem",
          borderTop: "1px solid #eee",
          paddingTop: "1rem",
        }}
      >
        <h2 style={{ fontSize: "1rem", marginBottom: "1rem" }}>
          コメント {comments.length}件
        </h2>

        {/* コメント投稿フォーム */}
        <form
          onSubmit={handleCommentSubmit}
          style={{ display: "flex", gap: 8, marginBottom: "1.5rem" }}
        >
          <input
            type="text"
            value={commentBody}
            onChange={(e) => setCommentBody(e.target.value)}
            placeholder="コメントを入力..."
            style={{
              flex: 1,
              padding: "8px 12px",
              border: "1px solid #ddd",
              borderRadius: 6,
              fontSize: "0.9rem",
            }}
          />
          <button
            type="submit"
            style={{
              padding: "8px 16px",
              background: "#e00",
              color: "#fff",
              border: "none",
              borderRadius: 6,
              cursor: "pointer",
              fontSize: "0.9rem",
            }}
          >
            投稿
          </button>
        </form>

        {/* コメント一覧 */}
        {comments.length === 0 ? (
          <p style={{ color: "#aaa", fontSize: "0.9rem" }}>
            まだコメントはありません
          </p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
            {comments.map((comment) => (
              <div key={comment.id} style={{ display: "flex", gap: 12 }}>
                <div
                  style={{
                    width: 36,
                    height: 36,
                    borderRadius: "50%",
                    background: "#ddd",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "0.85rem",
                    fontWeight: "bold",
                    flexShrink: 0,
                  }}
                >
                  {(comment.username || "?")[0].toUpperCase()}
                </div>
                <div>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: 8,
                      marginBottom: 2,
                    }}
                  >
                    <span style={{ fontSize: "0.85rem", fontWeight: "bold" }}>
                      {comment.username}
                    </span>
                    <span style={{ fontSize: "0.75rem", color: "#aaa" }}>
                      {new Date(comment.created_at).toLocaleDateString("ja-JP")}
                    </span>
                  </div>
                  <p style={{ margin: 0, fontSize: "0.9rem", lineHeight: 1.6 }}>
                    {comment.body}
                  </p>
                  <div style={{ display: "flex", gap: 6, marginTop: 6 }}>
                    <button
                      onClick={() => handleCommentLike(comment.id)}
                      style={{
                        background: "none",
                        border: `1px solid ${comment.liked ? "#e00" : "#ddd"}`,
                        borderRadius: 20,
                        cursor: "pointer",
                        color: comment.liked ? "#e00" : "#888",
                        fontSize: "0.8rem",
                        padding: "3px 10px",
                        display: "flex",
                        alignItems: "center",
                        gap: 4,
                      }}
                    >
                      👍 {comment.like_count}
                    </button>
                    <button
                      onClick={() => handleCommentDislike(comment.id)}
                      style={{
                        background: "none",
                        border: `1px solid ${comment.disliked ? "#555" : "#ddd"}`,
                        borderRadius: 20,
                        cursor: "pointer",
                        color: comment.disliked ? "#555" : "#888",
                        fontSize: "0.8rem",
                        padding: "3px 10px",
                        display: "flex",
                        alignItems: "center",
                        gap: 4,
                      }}
                    >
                      👎 {comment.dislike_count}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
