"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type UploadState = {
  label: string;
  percent: number | null; // null = 進捗不明のスピナー
};

// XHRを使ってMinIOに直接PUT（進捗取得のため）
function uploadWithProgress(
  presignedUrl: string,
  file: File,
  onProgress: (percent: number) => void
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", presignedUrl);
    xhr.setRequestHeader("Content-Type", file.type);

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        onProgress(Math.round((e.loaded / e.total) * 100));
      }
    };
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) resolve();
      else reject(new Error(`アップロードに失敗しました (${xhr.status})`));
    };
    xhr.onerror = () => reject(new Error("ネットワークエラーが発生しました"));
    xhr.send(file);
  });
}

export default function NewVideoPage() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [thumbnailFile, setThumbnailFile] = useState<File | null>(null);
  const [error, setError] = useState("");
  const [uploadState, setUploadState] = useState<UploadState | null>(null);

  const setStep = (label: string, percent: number | null = null) =>
    setUploadState({ label, percent });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const token = localStorage.getItem("token");
    if (!token) { router.push("/auth"); return; }

    try {
      // 1. メタデータ作成
      setStep("動画情報を保存中...");
      const createRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ title, description, thumbnail_url: "" }),
      });
      if (!createRes.ok) { setError(await createRes.text()); setUploadState(null); return; }
      const video = await createRes.json();

      // 2. サムネイルアップロード
      let thumbnailUrl = "";
      if (thumbnailFile) {
        setStep("サムネイルをアップロード中...", 0);
        const thumbRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/thumbnail-upload-url`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const { upload_url, thumbnail_url } = await thumbRes.json();
        await uploadWithProgress(upload_url, thumbnailFile, (p) =>
          setStep("サムネイルをアップロード中...", p)
        );
        thumbnailUrl = thumbnail_url;
      }

      // 3. 動画アップロード（進捗バーあり）
      let videoUrl = "";
      if (videoFile) {
        setStep("動画をアップロード中...", 0);
        const uploadRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/upload-url`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const { upload_url, video_url } = await uploadRes.json();
        await uploadWithProgress(upload_url, videoFile, (p) =>
          setStep("動画をアップロード中...", p)
        );
        videoUrl = video_url;
      }

      // 4. メタデータ更新
      if (thumbnailUrl || videoUrl) {
        setStep("情報を更新中...");
        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
          body: JSON.stringify({ title, description, thumbnail_url: thumbnailUrl || "", video_url: videoUrl || "" }),
        });
      }

      // 5. トランスコードをキューに投げる
      if (videoUrl) {
        setStep("変換キューに登録中...");
        await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${video.id}/transcode`, {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
        });
      }

      router.push(`/videos/${video.id}`);
    } catch (e) {
      setError(e instanceof Error ? e.message : "エラーが発生しました");
      setUploadState(null);
    }
  };

  const busy = uploadState !== null;
  const inputStyle = {
    width: "100%", padding: "0.5rem", fontSize: "1rem",
    boxSizing: "border-box" as const, borderRadius: 4,
    border: "1px solid #ddd",
  };

  return (
    <main style={{ maxWidth: 600, margin: "0 auto", padding: "2rem" }}>
      <h1 style={{ marginBottom: "1.5rem" }}>動画を投稿</h1>

      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }}>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem", fontWeight: 600, fontSize: "0.9rem" }}>
            タイトル <span style={{ color: "#e00" }}>*</span>
          </label>
          <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} required style={inputStyle} />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem", fontWeight: 600, fontSize: "0.9rem" }}>説明</label>
          <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={4} style={inputStyle} />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem", fontWeight: 600, fontSize: "0.9rem" }}>
            動画ファイル
          </label>
          <input type="file" accept="video/*" onChange={(e) => setVideoFile(e.target.files?.[0] ?? null)} style={inputStyle} />
          {videoFile && (
            <p style={{ margin: "0.3rem 0 0", fontSize: "0.8rem", color: "#888" }}>
              {videoFile.name}（{(videoFile.size / 1024 / 1024).toFixed(1)} MB）
            </p>
          )}
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem", fontWeight: 600, fontSize: "0.9rem" }}>
            サムネイル画像
            <span style={{ fontWeight: "normal", color: "#aaa", marginLeft: "0.5rem", fontSize: "0.8rem" }}>
              未選択の場合は自動生成
            </span>
          </label>
          <input type="file" accept="image/*" onChange={(e) => setThumbnailFile(e.target.files?.[0] ?? null)} style={inputStyle} />
        </div>

        {/* 進捗エリア */}
        {uploadState && (
          <div style={{ background: "#f5f5f5", borderRadius: 8, padding: "1rem" }}>
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "0.5rem" }}>
              <span style={{ fontSize: "0.85rem", color: "#555" }}>⏳ {uploadState.label}</span>
              {uploadState.percent !== null && (
                <span style={{ fontSize: "0.85rem", fontWeight: 600, color: "#333" }}>
                  {uploadState.percent}%
                </span>
              )}
            </div>
            {uploadState.percent !== null ? (
              /* 進捗バー */
              <div style={{ background: "#e0e0e0", borderRadius: 100, height: 8, overflow: "hidden" }}>
                <div style={{
                  height: "100%",
                  width: `${uploadState.percent}%`,
                  background: uploadState.percent === 100 ? "#22c55e" : "#e00",
                  borderRadius: 100,
                  transition: "width 0.2s ease",
                }} />
              </div>
            ) : (
              /* インジケーター（進捗不明時） */
              <div style={{ background: "#e0e0e0", borderRadius: 100, height: 8, overflow: "hidden" }}>
                <div style={{
                  height: "100%", width: "40%",
                  background: "#e00", borderRadius: 100,
                  animation: "slide 1.2s infinite ease-in-out",
                }} />
              </div>
            )}
          </div>
        )}

        <style>{`
          @keyframes slide {
            0%   { transform: translateX(-100%); }
            100% { transform: translateX(350%); }
          }
        `}</style>

        {error && <p style={{ color: "#e00", margin: 0, fontSize: "0.9rem" }}>{error}</p>}

        <button
          type="submit"
          disabled={busy}
          style={{
            padding: "0.8rem", fontSize: "1rem",
            background: busy ? "#ccc" : "#e00",
            color: "#fff", border: "none", borderRadius: 6,
            cursor: busy ? "not-allowed" : "pointer",
            fontWeight: 600,
          }}
        >
          {busy ? "投稿中..." : "投稿する"}
        </button>
      </form>
    </main>
  );
}
