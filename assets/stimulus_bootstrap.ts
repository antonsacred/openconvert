import { startStimulusApp } from '@symfony/stimulus-bundle';
import ConversionSelectorController from './controllers/conversion_selector_controller.ts';
import UploadQueueController from './controllers/upload_queue_controller.ts';

const app = startStimulusApp();
app.register('conversion-selector', ConversionSelectorController);
app.register('upload-queue', UploadQueueController);
