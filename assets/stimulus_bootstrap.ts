import { startStimulusApp } from '@symfony/stimulus-bundle';
import ConversionSelectorController from './controllers/conversion_selector_controller.ts';

const app = startStimulusApp();
app.register('conversion-selector', ConversionSelectorController);
