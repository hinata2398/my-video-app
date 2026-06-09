"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import HlsPlayer from "../../components/HlsPlayer";

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
  created_at: string;
};

export default function VideoDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [likeCount, setLikeCount] = useState(0); // いいね数
  const [liked, setLiked] = useState(false); // 自分がいいね済みか
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
    // いいね数を取得（誰でも見られる）
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}/likes`,
    );
    const data = await res.json();
    setLikeCount(data.count);

    // 自分がいいね済みか確認（ログイン中のみ）
    const token = localStorage.getItem("token");
    if (token) {
      // Toggle の代わりに Count エンドポイントは liked を返さないので
      // ここは liked の取得は省略してもOK（ボタン押した時に更新する）
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

    // APIから返ってきた最新の状態で更新
    setLiked(data.liked);
    setLikeCount(data.count);
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

  // useEffect でページ読み込み時に呼ぶ
  useEffect(() => {
    fetchVideo();
    fetchLikes();
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
        <button
          onClick={handleLike}
          style={{
            marginLeft: "auto",
            padding: "0.5rem 1rem",
            background: liked ? "#fff0f5" : "#fff",
            color: liked ? "#d63384" : "#888",
            border: "1px solid",
            borderColor: liked ? "#f9a8c9" : "#ddd",
            borderRadius: 20,
            cursor: "pointer",
            fontSize: "0.9rem",
            display: "flex",
            alignItems: "center",
            gap: "0.3rem",
          }}
        >
          ❤️ {likeCount}
        </button>
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
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
