module.exports = {
  presets: [
    '@babel/preset-env',
    [
      '@babel/preset-react',
      { runtime: 'automatic' }  // Enables the new JSX transform
    ],
    '@babel/preset-typescript'
  ],
};
