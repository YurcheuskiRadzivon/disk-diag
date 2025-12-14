// Global data storage for export
let cachedDisksData = null;
let cachedSmartData = null;
let cachedDiagnosticData = null;
let cachedBenchmarkData = null;

// Export menu toggle
function toggleExportMenu(tab) {
    const menu = document.getElementById(`export-menu-${tab}`);
    document.querySelectorAll('.dropdown-content').forEach(m => {
        if (m !== menu) m.classList.remove('show');
    });
    menu.classList.toggle('show');
}

// Close dropdowns when clicking outside
document.addEventListener('click', (e) => {
    if (!e.target.closest('.export-dropdown')) {
        document.querySelectorAll('.dropdown-content').forEach(m => m.classList.remove('show'));
    }
});

// Export data function
async function exportData(tab, format) {
    let data = null;
    let title = '';

    switch (tab) {
        case 'disks':
            data = cachedDisksData;
            title = 'Disks and Partitions Report';
            break;
        case 'smart':
            data = { smart: cachedSmartData, diagnostic: cachedDiagnosticData };
            title = 'SMART and Diagnostics Report';
            break;
        case 'benchmarks':
            data = cachedBenchmarkData;
            title = 'Benchmark Results Report';
            break;
    }

    if (!data) {
        alert('No data available for export. Please load the data first.');
        return;
    }

    try {
        const response = await fetch('/export/', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ data, title, format })
        });

        if (!response.ok) throw new Error('Export failed');

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `export.${format}`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);

    } catch (error) {
        alert(`Export error: ${error.message}`);
    }

    // Close menu
    document.querySelectorAll('.dropdown-content').forEach(m => m.classList.remove('show'));
}

document.addEventListener('DOMContentLoaded', () => {
    // Tab Switching Logic
    const tabs = document.querySelectorAll('.tab-btn');
    const contents = document.querySelectorAll('.tab-content');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => t.classList.remove('active'));
            contents.forEach(c => c.classList.remove('active'));

            tab.classList.add('active');
            const targetId = tab.getAttribute('data-tab');
            document.getElementById(targetId).classList.add('active');

            // Load data when tab is activated
            if (targetId === 'disks') loadDisks();
            if (targetId === 'smart') loadSmartInfo();
        });
    });

    // Initial Load
    loadDisks();

    // --- Tab 1: Disks & Partitions ---
    async function loadDisks() {
        const diskList = document.getElementById('disk-list');
        diskList.innerHTML = '<div class="loading">Loading disks...</div>';

        try {
            const response = await fetch('/base/disks');
            if (!response.ok) throw new Error('Failed to fetch disks');
            const disks = await response.json();

            // Cache for export
            cachedDisksData = { disks: disks, partitions: [] };

            diskList.innerHTML = '';

            if (!disks || disks.length === 0) {
                diskList.innerHTML = '<p>No disks found.</p>';
                return;
            }

            // Iterate through disks and fetch details for each
            for (let i = 0; i < disks.length; i++) {
                const disk = disks[i];
                // Note: Assuming the API returns a list of objects that have an ID or index. 
                // If not, we might need to use the index 'i'. 
                // The user prompt implies using /cdiskinfo/:id and /partitions/:id

                const diskItem = document.createElement('div');
                diskItem.className = 'disk-item';

                // Header
                const header = document.createElement('div');
                header.className = 'disk-header';
                header.innerHTML = `
                    <strong>Disk ${i}</strong>
                    <span>Click to expand</span>
                `;

                // Details Container
                const details = document.createElement('div');
                details.className = 'disk-details';
                details.id = `disk-details-${i}`;
                details.innerHTML = '<div class="loading">Loading details...</div>';

                header.addEventListener('click', () => {
                    const isOpen = details.classList.contains('open');
                    // Close all others (optional, but keeps UI clean)
                    // document.querySelectorAll('.disk-details').forEach(d => d.classList.remove('open'));

                    if (!isOpen) {
                        details.classList.add('open');
                        loadDiskDetails(i, details);
                    } else {
                        details.classList.remove('open');
                    }
                });

                diskItem.appendChild(header);
                diskItem.appendChild(details);
                diskList.appendChild(diskItem);
            }

        } catch (error) {
            console.error(error);
            diskList.innerHTML = `<p class="error">Error loading disks: ${error.message}</p>`;
        }
    }

    async function loadDiskDetails(id, container) {
        try {
            // Fetch Basic Info from multiple sources
            const [disksRes, cInfoRes, partRes] = await Promise.all([
                fetch('/base/disks'),
                fetch(`/base/cdiskinfo/${id}`),
                fetch(`/base/partitions/${id}`)
            ]);

            const disks = await disksRes.json();
            const cInfo = await cInfoRes.json();
            let partitions = [];
            if (partRes.ok) {
                partitions = await partRes.json();
            }

            // Find the specific disk info from the disks list
            const diskBasic = disks.find(d => d.index === id) || {};

            // Merge info
            const fullInfo = { ...diskBasic, ...cInfo };

            // Tooltips for disk info fields (Russian)
            const diskInfoTooltips = {
                index: "Внутренний индекс диска в системе",
                model: "Модель накопителя",
                serialNumber: "Серийный номер устройства",
                size: "Общий размер накопителя в байтах",
                sectorSize: "Размер сектора в байтах",
                mediaType: "Тип носителя (SSD/HDD)",
                busType: "Тип шины подключения (SATA/NVMe/USB)",
                firmwareRevision: "Версия прошивки накопителя",
                partitionStyle: "Стиль разметки (GPT/MBR)",
                healthStatus: "Общее состояние здоровья диска",
                operationalStatus: "Операционный статус устройства",
                temperature: "Текущая температура накопителя",
                controllerPath: "Путь контроллера в системе",
                name: "Название устройства",
                manufacturer: "Производитель",
                interfaceType: "Тип интерфейса подключения",
                totalCylinders: "Общее количество цилиндров",
                totalHeads: "Общее количество головок",
                totalSectors: "Общее количество секторов",
                totalTracks: "Общее количество дорожек",
                tracksPerCylinder: "Дорожек на цилиндр",
                sectorsPerTrack: "Секторов на дорожку",
                bytesPerSector: "Байт на сектор"
            };

            let html = `<h3>Basic Information</h3>`;
            html += `<table class="smart-table">`;
            for (const [key, value] of Object.entries(fullInfo)) {
                // Skip complex objects or large arrays if any
                if (typeof value === 'object' && value !== null) continue;
                const tooltip = diskInfoTooltips[key] || key;
                html += `<tr title="${tooltip}"><td><strong>${key}</strong></td><td>${value}</td></tr>`;
            }
            html += `</table>`;

            if (partitions && partitions.length > 0) {
                html += `<h3>Partitions</h3>`;
                html += `<div class="table-container">
                            <table class="partition-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Letter</th>
                                    <th>Label</th>
                                    <th>FS</th>
                                    <th>Size</th>
                                    <th>Free</th>
                                    <th>Type</th>
                                    <th>Boot</th>
                                    <th>Hidden</th>
                                </tr>
                            </thead>
                            <tbody>`;
                partitions.forEach(p => {
                    html += `<tr>
                                <td>${p.partitionId}</td>
                                <td>${p.driveLetter || '-'}</td>
                                <td>${p.label || '-'}</td>
                                <td>${p.fileSystem}</td>
                                <td>${formatBytes(p.size)}</td>
                                <td>${formatBytes(p.freeSpace)}</td>
                                <td>${p.type}</td>
                                <td>${p.isBoot ? 'Yes' : 'No'}</td>
                                <td>${p.isHidden ? 'Yes' : 'No'}</td>
                             </tr>`;
                });
                html += `</tbody></table></div>`;
            } else {
                html += `<p>No partitions found.</p>`;
            }

            container.innerHTML = html;

        } catch (error) {
            container.innerHTML = `<p class="error">Error loading details: ${error.message}</p>`;
        }
    }

    function formatBytes(bytes, decimals = 2) {
        if (!+bytes) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
    }


    // --- Tab 2: SMART & Diagnostics ---
    async function loadSmartInfo() {
        const container = document.getElementById('smart-data');
        const diskId = 0;

        try {
            const response = await fetch(`/smart/${diskId}`);
            if (!response.ok) throw new Error('Failed to fetch SMART data');
            const data = await response.json();

            // Cache for export
            cachedSmartData = data;

            // Display as formatted JSON
            container.innerHTML = `<pre class="json-view">${JSON.stringify(data, null, 2)}</pre>`;

        } catch (error) {
            container.innerHTML = `<p class="error">Error loading SMART data: ${error.message}</p>`;
        }
    }

    // Diagnostic Buttons
    document.getElementById('btn-diag-manual').addEventListener('click', () => runDiagnostic('manual'));
    document.getElementById('btn-diag-ai').addEventListener('click', () => runDiagnostic('gemini'));

    async function runDiagnostic(type) {
        const resultContainer = document.getElementById('diagnostic-result');
        resultContainer.innerHTML = '<div class="loading">Running diagnostic...</div>';
        const diskId = 0;

        const endpoint = type === 'gemini' ? `/diagnostic/gemini/${diskId}` : `/diagnostic/diagnostic/${diskId}`;

        try {
            const response = await fetch(endpoint);
            if (!response.ok) throw new Error('Diagnostic failed');
            const text = await response.text();

            let data;
            try {
                data = JSON.parse(text);
            } catch (e) {
                resultContainer.textContent = text;
                return;
            }

            // Status class mapping
            const statusLower = (data.status || 'unknown').toLowerCase();
            let statusClass = 'status-unknown';
            if (statusLower === 'ok' || statusLower === 'good' || statusLower === 'healthy') statusClass = 'status-ok';
            else if (statusLower === 'warning') statusClass = 'status-warning';
            else if (statusLower === 'critical' || statusLower === 'error' || statusLower === 'bad') statusClass = 'status-critical';

            // Problems HTML
            let problemsHtml = '';
            if (data.problems && data.problems.length > 0) {
                problemsHtml = `<div class="diag-section">
                    <div class="diag-section-title">Problems Found</div>
                    <ul class="diag-problems">${data.problems.map(p => `<li>${p}</li>`).join('')}</ul>
                </div>`;
            }

            // Metrics HTML with tooltips (Russian)
            const tooltips = {
                temperature_c: "Текущая температура накопителя в градусах Цельсия",
                life_remaining_percent: "Оставшийся ресурс SSD в процентах",
                data_written_tb: "Общий объём записанных данных в терабайтах",
                power_on_hours: "Общее время работы накопителя в часах",
                media_errors: "Количество неисправимых ошибок носителя",
                unsafe_shutdowns: "Количество небезопасных отключений питания"
            };

            let metricsHtml = '';
            if (data.metrics) {
                metricsHtml = `<div class="diag-section">
                    <div class="diag-section-title">Key Metrics</div>
                    <div class="diag-metrics">
                        ${Object.entries(data.metrics).map(([key, value]) => `
                            <div class="metric-card" title="${tooltips[key] || key.replace(/_/g, ' ')}">
                                <div class="metric-value">${value}</div>
                                <div class="metric-label">${key.replace(/_/g, ' ')}</div>
                            </div>
                        `).join('')}
                    </div>
                </div>`;
            }

            const html = `
                <div class="diag-card">
                    <div class="diag-header">
                        <span class="diag-method">${data.method || 'Diagnostic'}</span>
                        <span class="diag-status ${statusClass}">${data.status || 'Unknown'}</span>
                    </div>
                    
                    <div class="diag-score">
                        <span class="score-value">${data.health_score || 0}</span>
                        <span class="score-label">/ 100 Health Score</span>
                    </div>

                    <div class="diag-summary">
                        <div class="diag-section-title">Summary</div>
                        <p>${data.summary || 'No summary available.'}</p>
                    </div>

                    ${problemsHtml}
                    ${metricsHtml}
                </div>
            `;

            resultContainer.innerHTML = html;

            // Cache for export
            cachedDiagnosticData = data;

        } catch (error) {
            resultContainer.innerHTML = `<p class="error">Error: ${error.message}</p>`;
        }
    }


    // --- Tab 3: Benchmarks ---
    const benchForm = document.getElementById('benchmark-form');
    benchForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const type = document.getElementById('bench-type').value;
        const runs = document.getElementById('bench-runs').value;
        const seconds = document.getElementById('bench-seconds').value;
        const threads = document.getElementById('bench-threads').value;
        const size = document.getElementById('bench-size').value;
        const dir = document.getElementById('bench-dir').value;
        const balloon = document.getElementById('bench-balloon').checked;

        const diskId = 0;

        // Construct Query Params
        const params = new URLSearchParams({
            Runs: runs,
            Seconds: seconds,
            Threads: threads,
            SizeGiB: size,
            Dir: dir,
            Balloon: balloon,
            IODuration: seconds // Using same time limit for IOPS
        });

        const endpoint = `/benchmark/${type}/${diskId}?${params.toString()}`;

        // UI Updates
        document.getElementById('benchmark-progress').classList.remove('hidden');
        document.getElementById('benchmark-results').classList.add('hidden');
        const progressBar = document.querySelector('.progress-bar');
        progressBar.style.width = '0%';

        // Simulate progress (fake) since we don't have a websocket or progress endpoint
        let progress = 0;
        const interval = setInterval(() => {
            if (progress < 90) {
                progress += 5;
                progressBar.style.width = `${progress}%`;
            }
        }, 500);

        try {
            const response = await fetch(endpoint);
            if (!response.ok) throw new Error('Benchmark failed');

            const result = await response.json(); // Assuming JSON response

            // Cache for export
            cachedBenchmarkData = result;

            clearInterval(interval);
            progressBar.style.width = '100%';

            setTimeout(() => {
                document.getElementById('benchmark-progress').classList.add('hidden');
                document.getElementById('benchmark-results').classList.remove('hidden');
                document.getElementById('benchmark-output').textContent = JSON.stringify(result, null, 2);
            }, 500);

        } catch (error) {
            clearInterval(interval);
            document.getElementById('benchmark-progress').classList.add('hidden');
            alert(`Benchmark Error: ${error.message}`);
        }
    });
});
