import { defineConfig, presetUno, presetIcons } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': 'middle'
      }
    })
  ],
  theme: {
    colors: {
      primary: {
        DEFAULT: '#2D5A27',
        light: '#4A8A44',
        dark: '#1E3D1A'
      },
      accent: {
        DEFAULT: '#E8A838',
        light: '#F5C76A',
        dark: '#C68A20'
      },
      neutral: {
        50: '#FFFFFF',
        100: '#F8F9FA',
        200: '#F5F0E8',
        800: '#333333',
        900: '#1A1A1A'
      }
    }
  }
})
