// Evolver 运维模块 (src/ops/)
// 非飞书依赖的便携式工具集，用于 evolver 生命周期管理和维护。

module.exports = {
    lifecycle: require('./lifecycle'),
    skillsMonitor: require('./skills_monitor'),
    cleanup: require('./cleanup'),
    trigger: require('./trigger'),
    commentary: require('./commentary'),
    selfRepair: require('./self_repair'),
};
