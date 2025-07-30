import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
    base: '',
    plugins: [react()],
    server: {
        proxy: {
            '/refresh': 'http://localhost:8080',
            '/login': 'http://localhost:8080',
            '/logout': 'http://localhost:8080',
        }
    }
})
