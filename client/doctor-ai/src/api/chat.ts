const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const API_KEY = import.meta.env.VITE_GROQAPIKEY;

export const chatApi = {
  sendMessage: async (message: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/chat/completions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${API_KEY}`
        },
        body: JSON.stringify({
          model: 'mixtral-8x7b-32768',
          messages: [
            {
              role: 'system',
              content: 'You are a helpful AI medical assistant. Provide clear, concise, and accurate responses to medical queries.'
            },
            {
              role: 'user',
              content: message
            }
          ],
          temperature: 0.7,
          max_tokens: 1000
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(
          errorData?.error?.message || 
          `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return {
        response: data.choices[0].message.content
      };
    } catch (error) {
      console.error('Error sending chat message:', error);
      throw error;
    }
  },
};
