"use client";

import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { ChangeEvent, FormEvent, useState } from "react";
import { getImage } from "../api/api";

interface ImageDisplayProps {
  loading: boolean;
  imageUrl: string | null;
}

function ImageDisplay({ loading, imageUrl }: ImageDisplayProps) {
  if (loading) {
    return <div>Loading...</div>;
  }

  if (imageUrl) {
    return (
      <img src={imageUrl} alt="GitHub User Image" className="w-8/12 h-auto" />
    );
  }

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

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(
      `![GitHub persona](${apiUrl}/create?username=${username})`
    );
  };

  const handleSubmit = async (e: FormEvent) => {
    console.log("Username submitted:", username);
    e.preventDefault();
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const fullUrl = `${apiUrl}/create?username=${username}`;
    console.log("Calling API:", fullUrl);
    setLoading(true);
    setResultText(`![GitHub persona](${apiUrl}/create?username=${username})`);
    try {
      console.log("Starting API call...");
      const statusCode = await getImage(username);
      console.log("API response status:", statusCode);
      if (statusCode !== 200) {
        throw new Error(`Failed to fetch image, status code: ${statusCode}`);
      }
      console.log("API call successful");
      // 画像URLを設定
      console.log("Setting image URL:", fullUrl);
      setImageUrl(fullUrl);
    } catch (error) {
      console.error("API call failed:", error);
    } finally {
      setLoading(false);
      console.log("Loading set to false");
    }
    console.log(loading);
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
      <div className="flex flex-col items-center justify-center z-40">
        <ImageDisplay loading={loading} imageUrl={imageUrl} />
      </div>
    </form>
  );
}

export default Username;
