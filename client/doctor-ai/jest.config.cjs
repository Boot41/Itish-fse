module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  transform: {
    "^.+\\.(js|jsx|ts|tsx)$": "babel-jest"
  },
  transformIgnorePatterns: [
    "/node_modules/(?!(lucide-react)/)"
  ],
  extensionsToTreatAsEsm: ['.ts', '.tsx'],
  globals: {
    'ts-jest': {
      useESM: true,
    },
  },
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
};
