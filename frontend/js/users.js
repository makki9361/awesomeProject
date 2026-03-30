let currentSearchId = null;

document.addEventListener('DOMContentLoaded', () => {
    const userRole = localStorage.getItem('userRole');
    if (userRole !== 'admin') {
        window.location.href = '/rules.html';
        return;
    }

    loadUsers();

    const searchInput = document.getElementById('user-id-search');
    if (searchInput) {
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                searchUserById();
            }
        });
    }
});

async function loadUsers() {
    if (currentSearchId) {
        return;
    }

    const usersList = document.getElementById('users-list');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');

    try {
        if (loading) loading.style.display = 'flex';
        if (error) error.style.display = 'none';

        const users = await api.getUsersList();

        if (!usersList) return;

        if (users.length === 0) {
            usersList.innerHTML = '<div class="empty-state">Пользователи не найдены</div>';
            return;
        }

        displayUsers(users);

    } catch (err) {
        if (error) {
            error.style.display = 'flex';
            error.textContent = 'Ошибка загрузки пользователей: ' + err.message;
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

async function searchUserById() {
    const searchInput = document.getElementById('user-id-search');
    const userId = searchInput.value.trim();

    if (!userId) {
        showError('Введите ID пользователя');
        return;
    }

    const usersList = document.getElementById('users-list');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');

    try {
        if (loading) loading.style.display = 'flex';
        if (error) error.style.display = 'none';

        const user = await api.getUser(userId);
        currentSearchId = userId;

        if (!usersList) return;

        usersList.innerHTML = `
            <div class="user-card-detail">
                <div class="user-info-detail">
                    <div class="user-name-large">${escapeHtml(user.name)}</div>
                    <div class="user-role-badge-large">${getRoleName(user.role)}</div>
                    <div class="user-meta">
                        <span>ID: ${user.id}</span>
                        <span>Создан: ${formatDate(user.created_at)}</span>
                        <span>Обновлен: ${formatDate(user.updated_at)}</span>
                    </div>
                </div>
                <div class="user-actions-detail">
                    <button onclick="editUser(${user.id}, '${escapeHtml(user.name)}', '${user.role}')" class="btn btn-secondary">Редактировать</button>
                    <button onclick="deleteUser(${user.id})" class="btn btn-danger" ${user.id == localStorage.getItem('userId') ? 'disabled' : ''}>Удалить</button>
                </div>
            </div>
        `;

        const searchBar = document.querySelector('.search-bar');
        if (searchBar && !document.getElementById('reset-search-btn')) {
            const resetBtn = document.createElement('button');
            resetBtn.id = 'reset-search-btn';
            resetBtn.textContent = 'Показать всех пользователей';
            resetBtn.className = 'btn btn-secondary';
            resetBtn.style.marginTop = '10px';
            resetBtn.onclick = resetUserSearch;
            searchBar.appendChild(resetBtn);
        }

    } catch (err) {
        if (error) {
            error.style.display = 'flex';
            error.textContent = `Пользователь с ID ${userId} не найден: ${err.message}`;
        }
        usersList.innerHTML = `<div class="empty-state">Пользователь с ID ${userId} не найден</div>`;
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

function resetUserSearch() {
    currentSearchId = null;
    const searchInput = document.getElementById('user-id-search');
    if (searchInput) {
        searchInput.value = '';
    }

    const resetBtn = document.getElementById('reset-search-btn');
    if (resetBtn) {
        resetBtn.remove();
    }

    loadUsers();
}

function displayUsers(users) {
    const usersList = document.getElementById('users-list');
    if (!usersList) return;

    usersList.innerHTML = users.map(user => `
        <div class="user-card-detail">
            <div class="user-info-detail">
                <div class="user-name-large">${escapeHtml(user.name)}</div>
                <div class="d-flex gap-2">
                    <div class="user-role-badge-large">ID-${user.id}</div>
                    <div class="user-role-badge-large">${getRoleName(user.role)}</div>
                </div>
                <div class="user-meta">
                    <span>Создан: ${formatDate(user.created_at)}</span>
                    <span>Обновлен: ${formatDate(user.updated_at)}</span>
                </div>
            </div>
            <div class="user-actions-detail">
                <button onclick="editUser(${user.id}, '${escapeHtml(user.name)}', '${user.role}')" class="btn btn-secondary">Редактировать</button>
                <button onclick="deleteUser(${user.id})" class="btn btn-danger" ${user.id == localStorage.getItem('userId') ? 'disabled' : ''}><i class="fas fa-trash"></i></button>
            </div>
        </div>
    `).join('');
}

function showCreateUserModal() {
    document.getElementById('modal-title').textContent = 'Создание пользователя';
    document.getElementById('user-id').value = '';
    document.getElementById('user-name').value = '';
    document.getElementById('user-role-select').value = 'viewer';
    document.getElementById('user-modal').style.display = 'flex';
}

function editUser(id, name, role) {
    document.getElementById('modal-title').textContent = 'Редактирование пользователя';
    document.getElementById('user-id').value = id;
    document.getElementById('user-name').value = name;
    document.getElementById('user-role-select').value = role;
    document.getElementById('user-modal').style.display = 'flex';
}

function closeUserModal() {
    document.getElementById('user-modal').style.display = 'none';
}

async function saveUser(event) {
    event.preventDefault();

    const id = document.getElementById('user-id').value;
    const name = document.getElementById('user-name').value.trim();
    const role = document.getElementById('user-role-select').value;

    if (!name) {
        showError('Имя пользователя обязательно');
        return;
    }

    try {
        if (id) {
            await api.updateUser(id, { name, role });
            alert('Пользователь успешно обновлен');
        } else {
            await api.createUser({ name, role });
            alert('Пользователь успешно создан');
        }

        closeUserModal();

        if (currentSearchId) {
            resetUserSearch();
        } else {
            await loadUsers();
        }

    } catch (err) {
        alert('Ошибка при сохранении пользователя: ' + err.message);
    }
}

async function deleteUser(id) {
    const currentUserId = localStorage.getItem('userId');

    if (id == currentUserId) {
        alert('Вы не можете удалить самого себя');
        return;
    }

    if (!confirm('Вы уверены, что хотите удалить этого пользователя? Все его правила также будут удалены.')) {
        return;
    }

    try {
        await api.deleteUser(id);
        alert('Пользователь успешно удален');

        if (currentSearchId == id) {
            resetUserSearch();
        } else {
            await loadUsers();
        }
    } catch (err) {
        alert('Ошибка при удалении пользователя: ' + err.message);
    }
}