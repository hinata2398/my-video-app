import type { Metadata } from "next";
import Header from "./components/Header";

export const metadata: Metadata = {
  title: "My Video App",
  description: "動画投稿サイト ポートフォリオ",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body style={{ margin: 0, fontFamily: "sans-serif", background: "#fafafa" }}>
        <Header />
        {children}
      </body>
    </html>
  );
}
