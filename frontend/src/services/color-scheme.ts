export const colorScheme = (() => {
    const dE = document.documentElement;
    const dcs = 'data-bs-theme';
    const pcs = 'prefers-color-scheme';

    const init = (): void => {
        if (window.matchMedia) {
            const prefersDark = window.matchMedia(`(${pcs}: dark)`);

            prefersDark.addEventListener('change', (e: MediaQueryListEvent) => {
                dE.setAttribute(dcs, e.matches ? 'dark' : 'light');
                updateTheme();
            });
        }
    };

    const updateTheme = () => {
        updateSvgFilter();
    };

    function updateSvgFilter() {
        const reactLightDarkElements = document.querySelectorAll('.react-light-dark');
        reactLightDarkElements.forEach((element) => {
            // Use the cached elements
            (element as HTMLImageElement).style.filter = window.matchMedia(`(${pcs}: dark)`).matches ? 'brightness(0.7) contrast(1.1)' : 'brightness(1) contrast(1)';
        });
    }

    // Method to check the initialization state
    const getScheme = (): 'light' | 'dark' | null => {
        return dE.getAttribute(dcs) as 'light' | 'dark' | null;
    };

    // Expose methods (init, getScheme) to be used outside
    return {
        init,
        getScheme,
        updateTheme,
    };
})();

// Automatically invoke init on import
colorScheme.init();
