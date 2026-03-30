let currentRuleId = null;
let currentRule = null;

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    currentRuleId = urlParams.get('id');

    if (!currentRuleId) {
        showError('ID правила не указан');
        setTimeout(() => {
            window.location.href = '/rules.html';
        }, 2000);
        return;
    }

    loadRule();
});

async function loadRule() {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const ruleContent = document.getElementById('rule-content');

    try {
        if (loading) loading.style.display = 'flex';
        if (error) error.style.display = 'none';

        currentRule = await api.getRule(currentRuleId);

        if (!ruleContent) return;

        document.getElementById('rule-title').textContent = currentRule.title;
        document.getElementById('rule-status').textContent = getStatusName(currentRule.status);
        document.getElementById('rule-status').className = `rule-status status-${currentRule.status}`;
        document.getElementById('rule-version').textContent = currentRule.version;
        document.getElementById('rule-category').textContent = currentRule.category_name;
        document.getElementById('rule-created-at').textContent = formatDate(currentRule.created_at);
        document.getElementById('rule-updated-at').textContent = formatDate(currentRule.updated_at);
        document.getElementById('rule-content-text').innerHTML = escapeHtml(currentRule.content).replace(/\n/g, '<br>');

        const userRole = localStorage.getItem('userRole');
        const ruleActions = document.getElementById('rule-actions');
        const publishBtn = document.getElementById('publish-btn');
        const archiveBtn = document.getElementById('archive-btn');

        if (ruleActions && userRole === 'admin') {
            ruleActions.style.display = 'flex';

            if (currentRule.status === 'draft') {
                if (publishBtn) publishBtn.style.display = 'inline-block';
                if (archiveBtn) archiveBtn.style.display = 'none';
            } else if (currentRule.status === 'published') {
                if (publishBtn) publishBtn.style.display = 'none';
                if (archiveBtn) archiveBtn.style.display = 'inline-block';
            } else {
                if (publishBtn) publishBtn.style.display = 'none';
                if (archiveBtn) archiveBtn.style.display = 'none';
            }
        }

        ruleContent.style.display = 'block';

    } catch (err) {
        if (error) {
            error.style.display = 'block';
            error.textContent = 'Ошибка загрузки правила: ' + err.message;
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

function editCurrentRule() {
    if (currentRuleId) {
        window.location.href = `/rule-form.html?id=${currentRuleId}`;
    }
}

async function publishRule() {
    if (!confirm('Опубликовать правило? После публикации оно станет доступно всем пользователям.')) {
        return;
    }

    try {
        await api.publishRule(currentRuleId, { status: 'published' });
        alert('Правило успешно опубликовано');
        loadRule();
    } catch (err) {
        alert('Ошибка при публикации правила: ' + err.message);
    }
}

async function deleteCurrentRule() {
    if (!confirm('Вы уверены, что хотите удалить это правило? Это действие нельзя отменить.')) {
        return;
    }

    try {
        await api.deleteRule(currentRuleId);
        alert('Правило успешно удалено');
        window.location.href = '/rules.html';
    } catch (err) {
        alert('Ошибка при удалении правила: ' + err.message);
    }
}