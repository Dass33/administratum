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
            '/register': 'http://localhost:8080',
            '/update_column': 'http://localhost:8080',
            '/add_column': 'http://localhost:8080',
            '/delete_column': 'http://localhost:8080',
            '/get_sheet': 'http://localhost:8080',
            '/get_project': 'http://localhost:8080',
            '/create_project': 'http://localhost:8080',
            '/create_sheet': 'http://localhost:8080',
            '/rename_sheet': 'http://localhost:8080',
            '/delete_sheet': 'http://localhost:8080',
            '/rename_project': 'http://localhost:8080',
            '/delete_project': 'http://localhost:8080',
            '/add_share': 'http://localhost:8080',
            '/delete_row': 'http://localhost:8080',
            '/json': 'http://localhost:8080',
            '/change_game_url': 'http://localhost:8080',
        }
    }
})
