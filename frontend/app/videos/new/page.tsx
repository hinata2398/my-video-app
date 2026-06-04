"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function NewVideoPage() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null);
  const [error, setError] = useState("");
  const [progress, setProgress] = useState("");

  const uploadToMinio = async (presignedUrl: string, file: File) => {
    const res = await fetch(presignedUrl, {
      method: "PUT",
      body: file,
      headers: { "Content-Type": file.type },
    });
    if (!res.ok) throw new Error("アップロードに失敗しました");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const token = localStorage.getItem("token");
    if (!token) { router.push("/auth"); return; }

    try {
      // 1. 動画メタデータを作成
      setProgress("動画情報を保存中...");
      const createRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ title, description, thumbnail_url: "" }),
      });
      if (!createRes.ok) { setError(await createRes.text()); setProgress(""); return; }
      const video = await createRes.json();

      // 2. サムネイルアップロード（あれば）
      let thumbnailUrl = "";
      if (thumbnailFile) {
        setProgress("サムネイルをアップロード中...");
        const thumbRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/thumbnail-upload-url`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const { upload_url, thumbnail_url } = await thumbRes.json();
        await uploadToMinio(upload_url, thumbnailFile);
        thumbnailUrl = thumbnail_url;
      }

      // 3. 動画アップロード（あれば）
      let videoUrl = "";
      if (videoFile) {
        setProgress("動画をアップロード中...");
        const uploadRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/upload-url`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const { upload_url, video_url } = await uploadRes.json();
        await uploadToMinio(upload_url, videoFile);
        videoUrl = video_url;
      }

      // 4. メタデータを更新
      if (thumbnailUrl || videoUrl) {
        setProgress("情報を更新中...");
        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
          body: JSON.stringify({
            title,
            description,
            thumbnail_url: thumbnailUrl || "",
            video_url: videoUrl || "",
          }),
        });
      }

      // 5. 動画あり → トランスコードをキューに投げる（即座に返る）
      if (videoUrl) {
        setProgress("動画変換をキューに登録中...");
        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/transcode`, {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
        });
      }

      // 即座に詳細ページへ遷移（変換はバックグラウンドで継続）
      router.push(`/videos/${video.id}`);
    } catch (e) {
      setError(e instanceof Error ? e.message : "エラーが発生しました");
      setProgress("");
    }
  };

  const inputStyle = {
    width: "100%", padding: "0.5rem", fontSize: "1rem", boxSizing: "border-box" as const,
  };

  return (
    <main style={{ maxWidth: 600, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
      <h1>動画を投稿</h1>
      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }}>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>タイトル *</label>
          <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} required style={inputStyle} />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>説明</label>
          <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={4} style={inputStyle} />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>動画ファイル（mp4）</label>
          <input
            type="file"
            accept="video/mp4,video/*"
            onChange={(e) => setVideoFile(e.target.files?.[0] ?? null)}
            style={inputStyle}
          />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>サムネイル画像</label>
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setThumbnailFile(e.target.files?.[0] ?? null)}
            style={inputStyle}
          />
        </div>
        {error && <p style={{ color: "red", margin: 0 }}>{error}</p>}
        {progress && <p style={{ color: "#888", margin: 0 }}>⏳ {progress}</p>}
        <button
          type="submit"
          disabled={!!progress}
          style={{ padding: "0.75rem", fontSize: "1rem", background: progress ? "#ccc" : "#e00", color: "#fff", border: "none", borderRadius: 4, cursor: progress ? "not-allowed" : "pointer" }}
        >
          {progress ? "投稿中..." : "投稿する"}
        </button>
      </form>
    </main>
  );
}
