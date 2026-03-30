if (window.location.pathname === '/' || window.location.pathname === '/index.html') {
    loadUsers();
}

(function showAdminLinksImmediately() {
    const userRole = localStorage.getItem('userRole');
    if (userRole === 'admin') {
        const adminElements = document.querySelectorAll('.admin-only');
        adminElements.forEach(el => {
            el.style.display = 'flex';
        });
    }
})();

document.addEventListener('DOMContentLoaded', () => {
    const currentPath = window.location.pathname;

    if (currentPath !== '/' && currentPath !== '/index.html') {
        checkAuth();
        loadCurrentUser();
    }
});

function checkAuth() {
    const userId = localStorage.getItem('userId');
    if (!userId) {
        window.location.href = '/index.html';
        return false;
    }
    return true;
}

async function loadCurrentUser() {
    try {
        console.log('Loading current user...');
        const user = await api.getCurrentUser();
        console.log('Current user loaded:', user);

        const userNameSpan = document.getElementById('user-name');
        if (userNameSpan) {
            userNameSpan.textContent = user.name;
        }

        const userRoleSpan = document.getElementById('user-role');
        if (userRoleSpan) {
            userRoleSpan.textContent = getRoleName(user.role);
        }

        const isAdmin = user.role === 'admin';
        console.log('Is admin:', isAdmin);

        const adminElements = document.querySelectorAll('.admin-only');
        console.log('Admin elements found:', adminElements.length);

        adminElements.forEach(el => {
            console.log('Setting display for:', el.href || el.id);
            el.style.display = isAdmin ? 'flex' : 'none';
        });

        localStorage.setItem('userRole', user.role);

    } catch (err) {
        console.error('Failed to load user:', err);
        if (err.message.includes('401') || err.message.includes('403')) {
            logout();
        } else {
            const savedRole = localStorage.getItem('userRole');
            if (savedRole === 'admin') {
                const adminElements = document.querySelectorAll('.admin-only');
                adminElements.forEach(el => {
                    el.style.display = 'flex';
                });
            }
        }
    }
}

async function loadUsers() {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const usersList = document.getElementById('users-list');

    try {
        if (loading) loading.style.display = 'block';
        if (error) error.style.display = 'none';

        const users = await api.getUsers();

        if (!usersList) return;

        if (users.length === 0) {
            usersList.innerHTML = '<p class="empty-state">Нет пользователей в системе</p>';
            return;
        }

        usersList.innerHTML = users.map(user => `
            <div class="user-card" onclick="login(${user.id}, '${escapeHtml(user.name)}', '${user.role}')">
                <div class="user-name">${escapeHtml(user.name)}</div>
                <div class="user-role">${getRoleName(user.role)}</div>
            </div>
        `).join('');

    } catch (err) {
        if (error) {
            error.style.display = 'block';
            error.textContent = 'Ошибка загрузки пользователей: ' + err.message;
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

function login(userId, userName, userRole) {
    console.log('Login:', userId, userName, userRole);
    localStorage.setItem('userId', userId);
    localStorage.setItem('userName', userName);
    localStorage.setItem('userRole', userRole);

    window.location.href = '/rules.html';
}

function logout() {
    localStorage.removeItem('userId');
    localStorage.removeItem('userName');
    localStorage.removeItem('userRole');
    window.location.href = '/index.html';
}