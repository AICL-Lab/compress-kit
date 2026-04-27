<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="ck-footer">
    <div class="footer-container">
      <div class="footer-brand">
        <img src="/logo.svg" alt="CompressKit" class="footer-logo" />
        <p class="footer-tagline">
          Production-ready compression algorithms in C++17, Go, and Rust
        </p>
        <div class="footer-social">
          <a href="https://github.com/LessUp/compress-kit" target="_blank" rel="noopener" class="social-link" title="GitHub">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
          </a>
        </div>
      </div>
      
      <div class="footer-links">
        <div class="footer-column">
          <h4>Documentation</h4>
          <ul>
            <li><a href="/en/guide/getting-started">Getting Started</a></li>
            <li><a href="/en/guide/algorithms">Algorithms</a></li>
            <li><a href="/en/benchmarks/results">Benchmarks</a></li>
            <li><a href="/en/api/go">API Reference</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Resources</h4>
          <ul>
            <li><a href="https://github.com/LessUp/compress-kit" target="_blank" rel="noopener">GitHub Repository</a></li>
            <li><a href="https://github.com/LessUp/compress-kit/releases" target="_blank" rel="noopener">Releases</a></li>
            <li><a href="https://github.com/LessUp/compress-kit/issues" target="_blank" rel="noopener">Issues</a></li>
            <li><a href="https://github.com/LessUp/compress-kit/discussions" target="_blank" rel="noopener">Discussions</a></li>
          </ul>
        </div>
        
        <div class="footer-column">
          <h4>Community</h4>
          <ul>
            <li><a href="/en/guide/contributing">Contributing</a></li>
            <li><a href="https://github.com/LessUp/compress-kit/blob/master/CODE_OF_CONDUCT.md" target="_blank" rel="noopener">Code of Conduct</a></li>
            <li><a href="https://github.com/LessUp/compress-kit/blob/master/SECURITY.md" target="_blank" rel="noopener">Security</a></li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="footer-bottom">
      <p class="copyright">
        © {{ currentYear }} <strong>CompressKit</strong> by LessUp. 
        Released under the MIT License.
      </p>
      <p class="built-with">
        Built with ❤️ using <a href="https://vitepress.dev" target="_blank" rel="noopener">VitePress</a>
      </p>
    </div>
  </footer>
</template>

<style scoped>
.ck-footer {
  margin-top: 4rem;
  padding: 3rem 1.5rem 1.5rem;
  background: linear-gradient(180deg, var(--vp-c-bg) 0%, var(--vp-c-bg-alt) 100%);
  border-top: 1px solid var(--vp-c-divider);
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 4rem;
  padding-bottom: 2rem;
}

.footer-brand {
  display: flex;
  flex-direction: column;
}

.footer-logo {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
}

.footer-tagline {
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--vp-c-text-2);
  margin: 0 0 1rem;
}

.footer-social {
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 0.5rem;
  background: var(--vp-c-bg-elv);
  border: 1px solid var(--vp-c-divider);
  color: var(--vp-c-text-2);
  transition: all 0.2s ease;
}

.social-link:hover {
  background: var(--vp-c-brand-1);
  border-color: var(--vp-c-brand-1);
  color: white;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.footer-column h4 {
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--vp-c-text-1);
  margin: 0 0 1rem;
}

.footer-column ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.footer-column li {
  margin: 0.5rem 0;
}

.footer-column a {
  font-size: 0.9375rem;
  color: var(--vp-c-text-2);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-column a:hover {
  color: var(--vp-c-brand-1);
}

.footer-bottom {
  max-width: 1200px;
  margin: 0 auto;
  padding-top: 1.5rem;
  border-top: 1px solid var(--vp-c-divider);
  text-align: center;
}

.copyright {
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
  margin: 0 0 0.5rem;
}

.copyright strong {
  color: var(--vp-c-text-1);
}

.built-with {
  font-size: 0.8125rem;
  color: var(--vp-c-text-3);
  margin: 0;
}

.built-with a {
  color: var(--vp-c-text-2);
  text-decoration: none;
}

.built-with a:hover {
  color: var(--vp-c-brand-1);
}

@media (max-width: 960px) {
  .footer-container {
    grid-template-columns: 1fr;
    gap: 2rem;
  }
  
  .footer-links {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .footer-links {
    grid-template-columns: 1fr;
  }
  
  .ck-footer {
    padding: 2rem 1rem 1rem;
  }
}
</style>
