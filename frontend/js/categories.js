document.addEventListener('DOMContentLoaded', () => {
    loadCategories();
});

async function loadCategories() {
    const categoriesList = document.getElementById('categories-list');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');

    try {
        if (loading) loading.style.display = 'flex';
        if (error) error.style.display = 'none';

        const categories = await api.getCategories();

        if (!categoriesList) return;

        if (categories.length === 0) {
            categoriesList.innerHTML = '<div class="empty-state">Категории не найдены</div>';
            return;
        }

        categoriesList.innerHTML = categories.map(category => `
            <div class="category-card">
                <div class="category-name">${escapeHtml(category.name)}</div>
                <div class="category-actions">
                    <button onclick="editCategory(${category.id}, '${escapeHtml(category.name)}')" class="btn btn-secondary"><i class="fas fa-pen"></i></button>
                    <button onclick="deleteCategory(${category.id})" class="btn btn-danger"><i class="fas fa-trash"></i></button>
                </div>
            </div>
        `).join('');

    } catch (err) {
        if (error) {
            error.style.display = 'block';
            error.textContent = 'Ошибка загрузки категорий: ' + err.message;
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

function showCreateCategoryModal() {
    document.getElementById('modal-title').textContent = 'Создание категории';
    document.getElementById('category-id').value = '';
    document.getElementById('category-name').value = '';
    document.getElementById('category-modal').style.display = 'flex';
}

function editCategory(id, name) {
    document.getElementById('modal-title').textContent = 'Редактирование категории';
    document.getElementById('category-id').value = id;
    document.getElementById('category-name').value = name;
    document.getElementById('category-modal').style.display = 'flex';
}

function closeCategoryModal() {
    document.getElementById('category-modal').style.display = 'none';
}

async function saveCategory(event) {
    event.preventDefault();

    const id = document.getElementById('category-id').value;
    const name = document.getElementById('category-name').value.trim();

    if (!name) {
        showError('Название категории обязательно');
        return;
    }

    try {
        if (id) {
            await api.updateCategory(id, { name });
            alert('Категория успешно обновлена');
        } else {
            await api.createCategory({ name });
            alert('Категория успешно создана');
        }

        closeCategoryModal();
        await loadCategories();

    } catch (err) {
        alert('Ошибка при сохранении категории: ' + err.message);
    }
}

async function deleteCategory(id) {
    if (!confirm('Вы уверены, что хотите удалить эту категорию? Если в ней есть правила, удаление будет невозможно.')) {
        return;
    }

    try {
        await api.deleteCategory(id);
        alert('Категория успешно удалена');
        await loadCategories();
    } catch (err) {
        alert('Ошибка при удалении категории: ' + err.message);
    }
}