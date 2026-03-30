let ruleId = null;
let isEdit = false;

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    ruleId = urlParams.get('id');
    isEdit = !!ruleId;

    if (isEdit) {
        document.getElementById('page-title').textContent = 'Редактирование правила';
        document.getElementById('form-title').textContent = 'Редактирование правила';
        loadRule();
    } else {
        document.getElementById('page-title').textContent = 'Создание правила';
        document.getElementById('form-title').textContent = 'Создание правила';
    }

    loadCategories();
});

async function loadCategories() {
    try {
        const categories = await api.getCategories();
        const categorySelect = document.getElementById('category_id');

        if (categorySelect) {
            categorySelect.innerHTML = '<option value="">Выберите категорию</option>' +
                categories.map(cat => `<option value="${cat.id}">${escapeHtml(cat.name)}</option>`).join('');
        }
    } catch (err) {
        showError('Ошибка загрузки категорий: ' + err.message);
    }
}

async function loadRule() {
    try {
        showLoading(true);
        const rule = await api.getRule(ruleId);

        document.getElementById('title').value = rule.title;
        document.getElementById('category_id').value = rule.category_id;
        document.getElementById('status').value = rule.status;
        document.getElementById('content').value = rule.content;

    } catch (err) {
        showError('Ошибка загрузки правила: ' + err.message);
        setTimeout(() => {
            window.location.href = '/rules.html';
        }, 2000);
    } finally {
        showLoading(false);
    }
}

async function saveRule(event) {
    event.preventDefault();

    const title = document.getElementById('title').value.trim();
    const categoryId = document.getElementById('category_id').value;
    const status = document.getElementById('status').value;
    const content = document.getElementById('content').value.trim();

    if (!title) {
        showError('Заголовок обязателен для заполнения');
        return;
    }

    if (!categoryId) {
        showError('Выберите категорию');
        return;
    }

    if (!content) {
        showError('Содержание обязательно для заполнения');
        return;
    }

    const data = {
        title,
        category_id: parseInt(categoryId),
        status,
        content
    };

    if (!isEdit) {
        data.created_by = parseInt(localStorage.getItem('userId'));
    }

    try {
        showLoading(true);

        if (isEdit) {
            await api.updateRule(ruleId, data);
            alert('Правило успешно обновлено');
        } else {
            await api.createRule(data);
            alert('Правило успешно создано');
        }

        window.location.href = '/rules.html';

    } catch (err) {
        showError('Ошибка при сохранении: ' + err.message);
    } finally {
        showLoading(false);
    }
}

function cancel() {
    window.location.href = '/rules.html';
}