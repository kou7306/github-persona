export const getImage = async (username: string): Promise<number> => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const fullUrl = `${apiUrl}/create?username=${username}`;
  console.log("API URL:", fullUrl);
  console.log(
    "Environment variable NEXT_PUBLIC_API_URL:",
    process.env.NEXT_PUBLIC_API_URL
  );

  try {
    console.log("Sending fetch request to:", fullUrl);
    const response = await fetch(fullUrl, {
      method: "GET",
      headers: {
        Accept: "image/png",
      },
      // タイムアウトを延長（画像生成には時間がかかる）
      signal: AbortSignal.timeout(120000), // 30秒から120秒（2分）に延長
    });
    console.log("Response status:", response.status);
    console.log("Response headers:", response.headers);
    console.log("Response ok:", response.ok);
    return response.status;
  } catch (error) {
    console.error("API call error:", error);
    throw error;
  }
};
