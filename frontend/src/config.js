const isProd = import.meta.env.PROD

export const config = {
  apiBaseUrl: isProd 
    ? 'http://139.155.97.74:8000' 
    : 'http://localhost:8000'
}
