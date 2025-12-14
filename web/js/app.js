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

            let html = `<h3>Basic Information</h3>`;
            html += `<table class="smart-table">`;
            for (const [key, value] of Object.entries(fullInfo)) {
                // Skip complex objects or large arrays if any
                if (typeof value === 'object' && value !== null) continue;
                html += `<tr><td><strong>${key}</strong></td><td>${value}</td></tr>`;
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
            // The response might be text or JSON. 
            const text = await response.text();

            // Try to parse as JSON to format nicely, otherwise show text
            try {
                const json = JSON.parse(text);
                resultContainer.textContent = JSON.stringify(json, null, 2);
            } catch (e) {
                resultContainer.textContent = text;
            }

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
