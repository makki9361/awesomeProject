let currentPage = 1;
let currentFilters = {};

document.addEventListener('DOMContentLoaded', () => {
    loadCategories();
    loadRules();
    setupEventListeners();
});

function setupEventListeners() {
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                applyFilters();
            }
        });
    }
}

async function loadCategories() {
    try {
        const categories = await api.getCategories();
        const categoryFilter = document.getElementById('category-filter');

        if (categoryFilter) {
            categoryFilter.innerHTML = '<option value="">Все категории</option>' +
                categories.map(cat => `<option value="${cat.id}">${escapeHtml(cat.name)}</option>`).join('');
        }
    } catch (err) {
        console.error('Failed to load categories:', err);
    }
}

async function loadRules() {
    const rulesGrid = document.getElementById('rules-grid');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const pagination = document.getElementById('pagination');

    try {
        if (loading) loading.style.display = 'flex';
        if (error) error.style.display = 'none';

        const params = {
            page: currentPage,
            page_size: 6,
            ...currentFilters
        };

        const rules = await api.getRules(params);

        if (!rulesGrid) return;

        if (rules.length === 0) {
            rulesGrid.innerHTML = '<div class="empty-state">Правила не найдены</div>';
            if (pagination) pagination.style.display = 'none';
            return;
        }

        rulesGrid.innerHTML = rules.map(rule => `
            <div class="rule-card" onclick="viewRule(${rule.id})">
                <div class="rule-title">${escapeHtml(rule.title)}</div>
                <div class="rule-content">${escapeHtml(rule.content.substring(0, 150))}${rule.content.length > 150 ? '...' : ''}</div>
                <div>
                    <div class="rule-meta">
                        <span>Создано: ${formatDate(rule.created_at)}</span>
                        <span>Версия: ${rule.version}</span>
                        
                    </div>
                    <div class="rule-meta">
                        <span class="category">${escapeHtml(rule.category_name)}</span>
                        <span class="rule-status status-${rule.status}">${getStatusName(rule.status)}</span>
                    </div>
                </div>
            </div>
        `).join('');

        if (pagination) pagination.style.display = 'flex';
        updatePagination(rules.length);

    } catch (err) {
        if (error) {
            error.style.display = 'block';
            error.textContent = 'Ошибка загрузки правил: ' + err.message;
        }
    } finally {
        if (loading) loading.style.display = 'none';
    }
}

function viewRule(id) {
    window.location.href = `/rule-detail.html?id=${id}`;
}

function createRule() {
    window.location.href = '/rule-form.html';
}

function editRule(id) {
    window.location.href = `/rule-form.html?id=${id}`;
}

function applyFilters() {
    currentFilters = {
        category_id: document.getElementById('category-filter')?.value || undefined,
        status: document.getElementById('status-filter')?.value || undefined,
        search: document.getElementById('search-input')?.value || undefined
    };

    Object.keys(currentFilters).forEach(key => {
        if (!currentFilters[key]) delete currentFilters[key];
    });

    console.log('Applied filters:', currentFilters);

    currentPage = 1;
    loadRules();
}

function resetFilters() {
    const categoryFilter = document.getElementById('category-filter');
    const statusFilter = document.getElementById('status-filter');
    const searchInput = document.getElementById('search-input');

    if (categoryFilter) categoryFilter.value = '';
    if (statusFilter) statusFilter.value = '';
    if (searchInput) searchInput.value = '';

    currentFilters = {};
    currentPage = 1;
    loadRules();
}

function updatePagination(rulesCount) {
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    const pageInfo = document.getElementById('page-info');

    if (prevBtn) {
        prevBtn.disabled = currentPage === 1;
    }
    if (nextBtn) {
        nextBtn.disabled = rulesCount < 6;
    }
    if (pageInfo) {
        pageInfo.textContent = `Страница ${currentPage}`;
    }
}

function changePage(delta) {
    const newPage = currentPage + delta;
    if (newPage >= 1) {
        currentPage = newPage;
        loadRules();
    }
}