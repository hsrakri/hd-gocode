// Update system stats every 5 seconds
function updateSystemStats() {
    fetch('/api/system/stats')
        .then(response => response.json())
        .then(data => {
            document.getElementById('cpu-usage').textContent = `${data.CPUUsage.toFixed(1)}%`;
            document.getElementById('memory-usage').textContent = `${data.MemoryUsage.toFixed(1)}%`;
            
            // Update processes table
            const processesTable = document.getElementById('processes-table');
            processesTable.innerHTML = data.TopProcesses.map(process => `
                <tr>
                    <td>${process.PID}</td>
                    <td>${process.Command}</td>
                    <td>${process.CPU.toFixed(1)}%</td>
                    <td>${process.Memory.toFixed(1)}%</td>
                    <td>
                        <button class="btn btn-sm btn-primary btn-action" onclick="showProcessDetails(${process.PID})">
                            Details
                        </button>
                    </td>
                </tr>
            `).join('');
            
            // Update ports table
            const portsTable = document.getElementById('ports-table');
            portsTable.innerHTML = data.OpenPorts.map(port => `
                <tr>
                    <td>${port.Protocol}</td>
                    <td>${port.Local}</td>
                    <td>${port.Remote}</td>
                    <td>${port.State}</td>
                </tr>
            `).join('');
        })
        .catch(error => console.error('Error updating system stats:', error));
}

// Handle network diagnostics form submission
document.getElementById('network-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const target = document.getElementById('target-input').value;
    if (!target) return;

    // Show loading state
    const submitButton = this.querySelector('button[type="submit"]');
    const originalText = submitButton.innerHTML;
    submitButton.innerHTML = '<span class="loading"></span> Running...';
    submitButton.disabled = true;

    // Send request
    fetch('/api/network/diagnostics', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `target=${encodeURIComponent(target)}`
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('ping-output').textContent = data.PingResults;
        document.getElementById('mtr-output').textContent = data.MTRResults;
        document.getElementById('dns-output').textContent = data.DNSResults;
    })
    .catch(error => {
        console.error('Error running network diagnostics:', error);
        alert('Error running network diagnostics. Please try again.');
    })
    .finally(() => {
        // Restore button state
        submitButton.innerHTML = originalText;
        submitButton.disabled = false;
    });
});

// Show process details in modal
function showProcessDetails(pid) {
    fetch(`/api/process/${pid}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('ps-output').textContent = data.ps;
            document.getElementById('lsof-output').textContent = data.lsof;
            
            // Show modal
            const modal = new bootstrap.Modal(document.getElementById('processModal'));
            modal.show();
        })
        .catch(error => {
            console.error('Error fetching process details:', error);
            alert('Error fetching process details. Please try again.');
        });
}

// Initialize tooltips
document.addEventListener('DOMContentLoaded', function() {
    // Start system stats updates
    updateSystemStats();
    setInterval(updateSystemStats, 5000);
    
    // Initialize Bootstrap tooltips
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
}); 