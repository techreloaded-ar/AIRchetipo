const path = require('node:path');
const { ManifestGenerator } = require('./installers/lib/core/manifest-generator');

async function regenerateManifests() {
  const generator = new ManifestGenerator();
  const targetDir = process.argv[2] || 'z1';
  const airDir = path.join(process.cwd(), targetDir, 'air');

  // List of modules to include in manifests
  const selectedModules = ['aib', 'aim'];

  console.log('Regenerating manifests with relative paths...');
  console.log('Target directory:', airDir);

  try {
    const result = await generator.generateManifests(airDir, selectedModules, [], { ides: [] });
    console.log('✓ Manifests generated successfully:');
    console.log(`  - ${result.workflows} workflows`);
    console.log(`  - ${result.agents} agents`);
    console.log(`  - ${result.tasks} tasks`);
    console.log(`  - ${result.files} files in files-manifest.csv`);
  } catch (error) {
    console.error('Error generating manifests:', error);
    process.exit(1);
  }
}

regenerateManifests();
