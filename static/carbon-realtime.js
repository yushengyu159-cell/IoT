// Real-time CO2 and Sensor Data Functions
let sensorDataChart = null;
let realtimeCO2Interval = null;

// Load real-time CO2 data
async function loadRealtimeCO2() {
    try {
        const userEmail = localStorage.getItem('esg_login_email') || localStorage.getItem('esg_email') || '';
        const url = userEmail ? `/api/carbon/realtime?email=${encodeURIComponent(userEmail)}` : '/api/carbon/realtime';
        
        const response = await fetch(url);
        const result = await response.json();
        
        if (result.code === 200 && result.data) {
            const data = result.data;
            
            // Update real-time CO2 display
            const valueEl = document.getElementById('realtimeCO2Value');
            const timeEl = document.getElementById('realtimeCO2Time');
            const emissionEl = document.getElementById('realtimeCO2Emission');
            
            if (valueEl) valueEl.textContent = data.currentValue ? data.currentValue.toFixed(1) : '--';
            if (timeEl) {
                const date = new Date(parseInt(data.timestamp));
                timeEl.textContent = `Last updated: ${date.toLocaleString()}`;
            }
            if (emissionEl) {
                emissionEl.textContent = `Carbon Emission: ${data.co2Emission ? data.co2Emission.toFixed(2) : '--'} tCO2e`;
            }
            
            // Color coding based on CO2 level
            if (valueEl && data.currentValue) {
                if (data.currentValue < 450) {
                    valueEl.style.color = '#4caf50'; // Green - Good
                } else if (data.currentValue < 600) {
                    valueEl.style.color = '#ffc107'; // Yellow - Moderate
                } else {
                    valueEl.style.color = '#ff6b6b'; // Red - High
                }
            }
        }
    } catch (error) {
        console.error('Failed to fetch real-time CO2:', error);
    }
}

// Load sensor data chart for last N days
async function loadSensorChartData(days) {
    try {
        // Update active button
        document.getElementById('btn1Day').classList.remove('active');
        document.getElementById('btn7Days').classList.remove('active');
        if (days === 1) {
            document.getElementById('btn1Day').classList.add('active');
        } else {
            document.getElementById('btn7Days').classList.add('active');
        }
        
        const userEmail = localStorage.getItem('esg_login_email') || localStorage.getItem('esg_email') || '';
        const url = userEmail ? `/api/carbon/sensor-data?days=${days}&email=${encodeURIComponent(userEmail)}` : `/api/carbon/sensor-data?days=${days}`;
        
        const response = await fetch(url);
        const result = await response.json();
        
        if (result.code === 200 && result.data) {
            const data = result.data;
            renderSensorDataChart(data.dataPoints, days);
            
            // Update info
            const infoEl = document.getElementById('sensorDataInfo');
            if (infoEl) {
                infoEl.textContent = `${data.count} data points | From: ${data.fromDate} | To: ${data.toDate}`;
            }
        }
    } catch (error) {
        console.error('Failed to fetch sensor data:', error);
    }
}

// Render sensor data chart
function renderSensorDataChart(dataPoints, days) {
    const ctx = document.getElementById('sensorDataChart').getContext('2d');
    
    if (sensorDataChart) {
        sensorDataChart.destroy();
    }
    
    // Prepare data - limit points for performance
    let displayPoints = dataPoints;
    if (dataPoints.length > 100) {
        // Sample every Nth point if too many
        const step = Math.ceil(dataPoints.length / 100);
        displayPoints = dataPoints.filter((_, index) => index % step === 0);
    }
    
    const labels = displayPoints.map(dp => {
        if (days === 1) {
            return dp.time; // Show time for 1-day view
        } else {
            return dp.date; // Show date for multi-day view
        }
    });
    const values = displayPoints.map(dp => dp.value);
    
    sensorDataChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'CO2 Level (ppm)',
                data: values,
                borderColor: '#2c7a62',
                backgroundColor: 'rgba(44, 122, 98, 0.1)',
                fill: true,
                tension: 0.4,
                pointRadius: days === 1 ? 3 : 2,
                pointHoverRadius: 6
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    callbacks: {
                        label: function(context) {
                            const idx = context.dataIndex;
                            const dp = displayPoints[idx];
                            return [`CO2: ${dp.value} ppm`, `Time: ${dp.timestamp}`, `Date: ${dp.date}`];
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: false,
                    title: {
                        display: true,
                        text: 'CO2 (ppm)'
                    },
                    grid: {
                        borderDash: [2, 4]
                    }
                },
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        maxTicksLimit: 12
                    }
                }
            },
            interaction: {
                mode: 'nearest',
                axis: 'x',
                intersect: false
            }
        }
    });
}

// Initialize real-time CO2 updates
function initRealtimeCO2Updates() {
    // Load immediately
    loadRealtimeCO2();
    loadSensorChartData(1);
    
    // Update every 30 seconds
    if (realtimeCO2Interval) {
        clearInterval(realtimeCO2Interval);
    }
    realtimeCO2Interval = setInterval(() => {
        loadRealtimeCO2();
    }, 30000);
}

// Stop real-time updates when leaving carbon transparency view
function stopRealtimeCO2Updates() {
    if (realtimeCO2Interval) {
        clearInterval(realtimeCO2Interval);
        realtimeCO2Interval = null;
    }
}

// Call this when switching to carbon transparency view
function onCarbonTransparencyViewActivated() {
    initRealtimeCO2Updates();
}

// Call this when switching away from carbon transparency view
function onCarbonTransparencyViewDeactivated() {
    stopRealtimeCO2Updates();
}
