"use client";

import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { ChangeEvent, FormEvent, useState, useEffect } from "react";
import { getImage } from "../api/api";

interface ImageDisplayProps {
  loading: boolean;
  imageUrl: string | null;
  onImageLoad: () => void;
  onImageError: () => void;
}

function ImageDisplay({
  loading,
  imageUrl,
  onImageLoad,
  onImageError,
}: ImageDisplayProps) {
  console.log("ImageDisplay - loading:", loading, "imageUrl:", imageUrl);

  if (loading) {
    console.log("Showing loading spinner");
    return (
      <div
        className="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-[9999]"
        style={{ position: "fixed", zIndex: 9999 }}
      >
        <div className="flex flex-col items-center justify-center p-12 bg-gray-900 bg-opacity-95 rounded-xl min-h-[300px] border-4 border-green-500 shadow-2xl">
          <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-green-500 mb-6"></div>
          <p className="text-white text-xl font-bold mb-2">画像を生成中...</p>
          <p className="text-gray-300 text-lg">しばらくお待ちください</p>
          <div className="mt-4 text-green-400 text-sm">処理中...</div>
        </div>
      </div>
    );
  }

  if (imageUrl) {
    console.log("Showing image with URL:", imageUrl);
    return (
      <img
        src={imageUrl}
        alt="GitHub User Image"
        className="w-8/12 h-auto"
        onLoad={(e) => {
          console.log("Image loaded successfully");
          const img = e.target as HTMLImageElement;
          console.log(
            "Image dimensions:",
            img.naturalWidth,
            "x",
            img.naturalHeight
          );
          // 画像が正常に読み込まれたことを通知
          onImageLoad();
        }}
        onError={(e) => {
          console.error("Image load error:", e);
          const img = e.target as HTMLImageElement;
          console.error("Failed image URL:", img.src);
          // 画像の読み込みに失敗したことを通知
          onImageError();
        }}
      />
    );
  }

  console.log("Showing nothing");
  return null;
}

function Username() {
  const [username, setUsername] = useState<string>("");
  const [imageUrl, setImageUrl] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const [resultText, setResultText] = useState<string>(
    `![GitHub persona](${apiUrl}/create?username=`
  );

  // ローディング状態を監視
  useEffect(() => {
    console.log("Loading state changed:", loading);
  }, [loading]);

  // 画像が読み込まれた時の処理
  const handleImageLoad = () => {
    console.log("Image loaded, setting loading to false");
    setLoading(false);
  };

  // 画像の読み込みに失敗した時の処理
  const handleImageError = () => {
    console.log("Image failed to load, setting loading to false");
    setLoading(false);
  };

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(
      `![GitHub persona](${apiUrl}/create?username=${username})`
    );
  };

  const handleSubmit = async (e: FormEvent) => {
    console.log("Username submitted:", username);
    e.preventDefault();
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const timestamp = Date.now(); // キャッシュバスティング用
    const fullUrl = `${apiUrl}/create?username=${username}&t=${timestamp}`;
    console.log("Calling API:", fullUrl);
    setLoading(true);
    console.log("Loading set to true");
    setImageUrl(""); // 古い画像をクリア
    setResultText(`![GitHub persona](${apiUrl}/create?username=${username})`);

    try {
      console.log("Starting API call...");
      const statusCode = await getImage(username);
      console.log("API response status:", statusCode);

      if (statusCode !== 200) {
        throw new Error(`Failed to fetch image, status code: ${statusCode}`);
      }

      console.log("API call successful, setting image URL");
      // API呼び出しが成功した後に画像URLを設定
      setImageUrl(fullUrl);

      // ローディングは画像のonLoadイベントで解除されるため、ここでは解除しない
      console.log("Image URL set, waiting for image to load...");
    } catch (error) {
      console.error("API call failed:", error);
      setLoading(false);
    }
  };

  return (
    <form
      className="w-auto flex flex-col items-center justify-center mb-4 space-y-3 text-black"
      onSubmit={handleSubmit}
    >
      <input
        value={username}
        type="text"
        onChange={(e: ChangeEvent<HTMLInputElement>) =>
          setUsername(e.target.value)
        }
        className="w-64 px-4 py-2 border rounded-lg focus:outline-none focus:border-green-400"
        placeholder="Username"
      />
      <button className="w-64 px-4 py-2 text-white bg-green-500 rounded transform transition-transform duration-200 hover:bg-green-400 hover:scale-95">
        Generate
      </button>
      <div className="App">
        {imageUrl && (
          <div className="relative bg-gray-800 p-6 rounded-md">
            <div className="absolute top-1 right-1">
              <Tooltip title="Copy to Clipboard" placement="top" arrow>
                <IconButton
                  color="primary"
                  size="small"
                  onClick={copyToClipboard}
                >
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </div>
            <p className="text-white w-72 h-auto px-4 resize-none bg-transparent border-none focus:outline-none">
              {resultText}
            </p>
          </div>
        )}
      </div>
      <div className="flex flex-col items-center justify-center z-[9999] relative min-h-[300px] w-full">
        <ImageDisplay
          loading={loading}
          imageUrl={imageUrl}
          onImageLoad={handleImageLoad}
          onImageError={handleImageError}
        />
      </div>
    </form>
  );
}

export default Username;
