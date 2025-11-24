// Retro Terminal Effects

// Typing animation for the terminal prompt
function typeText(element, text, speed = 50) {
    let i = 0;
    element.textContent = '';
    
    function type() {
        if (i < text.length) {
            element.textContent += text.charAt(i);
            i++;
            setTimeout(type, speed);
        }
    }
    
    type();
}

// Initialize typing animation when page loads
document.addEventListener('DOMContentLoaded', () => {
    const typingElement = document.querySelector('.typing');
    if (typingElement) {
        const text = typingElement.textContent;
        setTimeout(() => {
            typeText(typingElement, text, 80);
        }, 500);
    }

    // Add observer for fade-in animations
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }
        });
    }, observerOptions);

    // Observe all feature boxes and examples
    document.querySelectorAll('.feature-box, .example, .step').forEach(el => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'opacity 0.6s, transform 0.6s';
        observer.observe(el);
    });

    // Random CRT flicker effect
    setInterval(() => {
        if (Math.random() > 0.97) {
            document.body.style.opacity = '0.95';
            setTimeout(() => {
                document.body.style.opacity = '1';
            }, 50);
        }
    }, 100);

    // Add click effect to retro buttons
    document.querySelectorAll('.retro-button, .retro-link').forEach(button => {
        button.addEventListener('click', function(e) {
            this.style.transform = 'scale(0.95)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 100);
        });
    });

    // Easter egg: Konami code
    let konamiCode = [];
    const konamiSequence = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a'];
    
    document.addEventListener('keydown', (e) => {
        konamiCode.push(e.key);
        konamiCode = konamiCode.slice(-10);
        
        if (konamiCode.join(',') === konamiSequence.join(',')) {
            activateMatrixMode();
        }
    });
});

// Matrix rain easter egg
function activateMatrixMode() {
    const canvas = document.createElement('canvas');
    canvas.style.position = 'fixed';
    canvas.style.top = '0';
    canvas.style.left = '0';
    canvas.style.width = '100%';
    canvas.style.height = '100%';
    canvas.style.pointerEvents = 'none';
    canvas.style.zIndex = '9999';
    canvas.style.opacity = '0.3';
    document.body.appendChild(canvas);

    const ctx = canvas.getContext('2d');
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    const chars = '01';
    const fontSize = 14;
    const columns = canvas.width / fontSize;
    const drops = Array(Math.floor(columns)).fill(1);

    function draw() {
        ctx.fillStyle = 'rgba(0, 0, 0, 0.05)';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        
        ctx.fillStyle = '#00ff41';
        ctx.font = fontSize + 'px monospace';

        for (let i = 0; i < drops.length; i++) {
            const text = chars[Math.floor(Math.random() * chars.length)];
            ctx.fillText(text, i * fontSize, drops[i] * fontSize);
            
            if (drops[i] * fontSize > canvas.height && Math.random() > 0.975) {
                drops[i] = 0;
            }
            drops[i]++;
        }
    }

    const interval = setInterval(draw, 33);
    
    // Stop after 10 seconds
    setTimeout(() => {
        clearInterval(interval);
        canvas.remove();
    }, 10000);
}

// Handle window resize
window.addEventListener('resize', () => {
    // Update any dynamic sizing if needed
});

// Smooth scroll for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Add glitch effect on logo hover
const logo = document.querySelector('.ascii-logo');
if (logo) {
    logo.addEventListener('mouseenter', () => {
        logo.style.animation = 'glitch 0.3s infinite';
    });
    
    logo.addEventListener('mouseleave', () => {
        logo.style.animation = 'none';
    });
}

// Add glitch keyframes dynamically
const style = document.createElement('style');
style.textContent = `
    @keyframes glitch {
        0% {
            text-shadow: 0 0 10px var(--terminal-glow);
        }
        25% {
            text-shadow: -2px 0 10px var(--terminal-glow), 2px 0 10px var(--magenta);
        }
        50% {
            text-shadow: 2px 0 10px var(--terminal-glow), -2px 0 10px var(--cyan);
        }
        75% {
            text-shadow: -2px 0 10px var(--cyan), 2px 0 10px var(--amber);
        }
        100% {
            text-shadow: 0 0 10px var(--terminal-glow);
        }
    }
`;
document.head.appendChild(style);

console.log('%c┏━━━━━━━━━━━━━━━━━━━━━━━━━━┓', 'color: #00ff41; font-size: 14px; font-family: monospace;');
console.log('%c┃  WEATHERORNOT v1.0.0    ┃', 'color: #00ff41; font-size: 14px; font-family: monospace;');
console.log('%c┃  Retro Terminal Weather ┃', 'color: #00ff41; font-size: 14px; font-family: monospace;');
console.log('%c┃  Try the Konami code!   ┃', 'color: #ffb000; font-size: 14px; font-family: monospace;');
console.log('%c┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛', 'color: #00ff41; font-size: 14px; font-family: monospace;');

