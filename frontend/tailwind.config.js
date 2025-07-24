export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsz,tsx}"
    ],
    safelist: [
    ],
    theme: {
        extend: {
            colors: {
                'figma-black': '#0B1F42',
                'figma-stone': '#245375',
                'figma-light-gray': '#A7BAC8',
                'figma-gray': '#D9D9D9',
                'figma-honey': '#FFD32A',
                'figma-carrots': '#FF6130',
                'figma-berries': '#E1015B',
                'figma-rose': '#EB4C79',
                'figma-lavender': '#AA4BB3',
                'figma-pool': '#199CF9',
                'figma-teal': '#02A3A4',
                'figma-lime': '#0CB43F',
                'figma-forest': '#028661',
                'figma-winter': '#9FD7FF',
                'figma-indigo-40': '#ACABD3',
                'figma-indigo': '#2F2D91',
                'figma-pale': '#FF9CC4',
                'figma-white': '#EEF1F4',
                'figma-stone-40': '#A7BAC8',
                'figma-pool-40': '#A3D7FD',
                'figma-lime-40': '#9EE1B2',
                'figma-green': '#609C5D',

            },
            backgroundImage: {
                'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
            },
            screens: {
                'xs': '380px',
            },
            keyframes: {
                growAndShrink: {
                    "0%": { transform: "scale(1)" },
                    "50%": { transform: "scale(1.3)" },
                    "100%": { transform: "scale(1)" },
                },
            },
            animation: {
                "grow-shrink": "growAndShrink 0.3s ease-in-out",
            },
        },
    },
    plugins: [],
}

