// 快速诊断localStorage数据
console.log('🔍 ESG VISA 数据诊断开始...');

const requiredFields = [
    'esg_email',
    'esg_password', 
    'esg_fullName',
    'esg_role',
    'esg_building_name',
    'esg_building_addr',
    'esg_btype'
];

console.log('📊 当前localStorage数据状态:');
requiredFields.forEach(field => {
    const value = localStorage.getItem(field);
    if (value) {
        console.log(`✅ ${field}: ${field.includes('password') ? '***' : value}`);
    } else {
        console.log(`❌ ${field}: 缺失`);
    }
});

// 检查是否有数据
const hasData = requiredFields.every(field => localStorage.getItem(field));
console.log(`\n🎯 数据完整性: ${hasData ? '完整' : '不完整'}`);

if (!hasData) {
    console.log('\n🔧 建议操作:');
    console.log('1. 访问 /static/debug-data.html 查看详细状态');
    console.log('2. 返回 /static/index.html 重新开始');
    console.log('3. 确保每个步骤都正确填写并保存');
}
