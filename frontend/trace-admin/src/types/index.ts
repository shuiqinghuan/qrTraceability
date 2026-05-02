export interface Product {
  id: number
  name: string
  code: string
  planting_location: string
  planting_date: string
  images: string[]
  video: string | null
  media_files?: MediaFile[]
  harvest_start_date: string | null
  harvest_end_date: string | null
  sugar_content: number | null
  weight: number | null
  taste: string
  quality: string
  quality_summary: string
  suitable_for: string[]
  created_at: string
  updated_at: string
}

export interface MediaFile {
  id: number
  url: string
  media_type: 'Image' | 'Video'
  filename: string
  file_size: number
  uploaded_at: string
}

export interface LoginDTO {
  username: string
  password: string
}

export interface TokenResponse {
  access: string
  refresh: string
}
