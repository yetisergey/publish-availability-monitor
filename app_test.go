package main

import (
	"os"
	"testing"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/content"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	os.Exit(m.Run())
}

func TestIsMessagePastPublishSLA_pastSLA(t *testing.T) {
	publishDate := time.Now().Add(-(threshold + 1) * time.Second)
	if !isMessagePastPublishSLA(publishDate, threshold) {
		t.Error("Did not detect message past SLA")
	}
}

func TestIsMessagePastPublishSLA_notPastSLA(t *testing.T) {
	publishDate := time.Now()
	if isMessagePastPublishSLA(publishDate, threshold) {
		t.Error("Valid message marked as passed SLA")
	}
}

func TestIsIgnorableMessage_naturalMessage(t *testing.T) {
	if isIgnorableMessage(naturalTID) {
		t.Error("Normal message marked as ignorable")
	}
}

func TestIsIgnorableMessage_syntheticMessage(t *testing.T) {
	if !isIgnorableMessage(syntheticTID) {
		t.Error("Synthetic message marked as normal")
	}
}

func TestGetCredentials(t *testing.T) {
	environments["env1"] = Environment{"env1", "http://env1.example.org", "http://s3.example.org", "user1", "pass1"}
	environments["env2"] = Environment{"env2", "http://env2.example.org", "http://s3.example.org", "user2", "pass2"}

	username, password := getCredentials("http://env2.example.org/__some-service")
	if username != "user2" || password != "pass2" {
		t.Error("incorrect credentials returned")
	}
}

func TestGetCredentials_Unauthenticated(t *testing.T) {
	environments["env1"] = Environment{"env1", "http://env1.example.org", "http://s3.example.org", "user1", "pass1"}
	environments["env2"] = Environment{"env2", "http://env2.example.org", "http://s3.example.org", "user2", "pass2"}

	username, password := getCredentials("http://env3.example.org/__some-service")
	if username != "" || password != "" {
		t.Error("incorrect credentials returned")
	}
}

func TestGetValidationEndpointKey_CompoundStory(t *testing.T) {
	validationEndpointKey := getValidationEndpointKey(supportedSourceCodeCompoundStory, naturalTID, testUuid)
	assert.Equal(t, validationEndpointKey, "EOM::CompoundStory", "Didn't get expected validation url key for compound story")
}

func TestGetValidationEndpointKey_ContentPlaceholderCompoundStory(t *testing.T) {
	validationEndpointKey := getValidationEndpointKey(contentplaceHolderCompoundStory, naturalTID, testUuid)
	assert.Equal(t, validationEndpointKey, "EOM::CompoundStory_ContentPlaceholder", "Didn't get expected validation url key for content placeholder")
}

func TestGetValidationEndpointKey_Story(t *testing.T) {
	validationEndpointKey := getValidationEndpointKey(supportedSourceCodeStory, naturalTID, testUuid)
	assert.Equal(t, validationEndpointKey, "EOM::Story", "Didn't get expected validation url key for Stroy")
}

const threshold = 120
const syntheticTID = "SYNTHETIC-REQ-MONe4d2885f-1140-400b-9407-921e1c7378cd"
const naturalTID = "tid_xltcnbckvq"
const testUuid = "uuid"

var supportedSourceCodeCompoundStory = content.EomFile{
	UUID:             testUuid,
	Type:             "EOM::CompoundStory",
	ContentType:      "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       supportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}

var contentplaceHolderCompoundStory = content.EomFile{
	UUID:             testUuid,
	Type:             "EOM::CompoundStory_ContentPlaceholder",
	ContentType:      "EOM::CompoundStory",
	Value:            "bar",
	Attributes:       supportedSourceCodeAttributesContentPlaceholder,
	SystemAttributes: "systemAttributes",
}

var supportedSourceCodeStory = content.EomFile{
	UUID:             testUuid,
	Type:             "EOM::Story",
	ContentType:      "EOM::CompoundStory",
	Value:            "value",
	Attributes:       supportedSourceCodeAttributes,
	SystemAttributes: "systemAttributes",
}

var supportedSourceCodeAttributes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DILeadCompanies /><DITemporaryCompanies><DITemporaryCompany><DICoTempCode /><DICoTempDescriptor /><DICoTickerCode /></DITemporaryCompany></DITemporaryCompanies><DIFTSEGlobalClassifications /><DIStockExchangeIndices /><DIHotTopics /><DIHeadlineCopy>Global oil inventory stands at record level</DIHeadlineCopy><DIBylineCopy>Anjli Raval, Oil and Gas Correspondent</DIBylineCopy><DIFTNPSections /><DIFirstParCopy>International Energy Agency says crude stockpiles near 3bn barrels despite robust demand growth</DIFirstParCopy><DIMasterImgFileRef>/FT/Graphics/Online/Master_2048x1152/2015/08/New%20story%20of%20MAS_crude-oil_02.jpg?uuid=c61fcc32-61cd-11e5-9846-de406ccb37f2</DIMasterImgFileRef></EditorialDisplayIndexing><OutputChannels><DIFTN><DIFTNPublicationDate /><DIFTNZoneEdition /><DIFTNPage /><DIFTNTimeEdition /><DIFTNFronts /></DIFTN><DIFTcom><DIFTcomWebType>story</DIFTcomWebType><DIFTcomDisplayCodes><DIFTcomDisplayCodeRank1><DIFTcomDisplayCode title=\"Markets - Commodities\"><DIFTcomDisplayCodeFTCode>MKCO</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Commodities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Commodities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank1><DIFTcomDisplayCodeRank2><DIFTcomDisplayCode title=\"Markets - European Equities\"><DIFTcomDisplayCodeFTCode>MKEU</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - European Equities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>European Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy - Oil &amp; Gas\"><DIFTcomDisplayCodeFTCode>OG8E</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy - Oil &amp; Gas</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Oil &amp; Gas</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy\"><DIFTcomDisplayCodeFTCode>NDEM</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Energy</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Companies\"><DIFTcomDisplayCodeFTCode>BNIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Companies</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets\"><DIFTcomDisplayCodeFTCode>MKIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets - Equities Main\"><DIFTcomDisplayCodeFTCode>MKEQ</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Equities Main</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank2></DIFTcomDisplayCodes><DIFTcomSubscriptionLevel>0</DIFTcomSubscriptionLevel><DIFTcomUpdateTimeStamp>False</DIFTcomUpdateTimeStamp><DIFTcomIndexAndSynd>false</DIFTcomIndexAndSynd><DIFTcomSafeToSyndicate>True</DIFTcomSafeToSyndicate><DIFTcomInitialPublication>20151113105031</DIFTcomInitialPublication><DIFTcomLastPublication>20151113132953</DIFTcomLastPublication><DIFTcomSuppresInlineAds>False</DIFTcomSuppresInlineAds><DIFTcomMap>True</DIFTcomMap><DIFTcomDisplayStyle>Normal</DIFTcomDisplayStyle><DIFTcomFeatureType>Normal</DIFTcomFeatureType><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted><DIFTcomMakeUnlinkable>False</DIFTcomMakeUnlinkable><isBestStory>0</isBestStory><DIFTcomCMRId>3040370</DIFTcomCMRId><DIFTcomCMRHint /><DIFTcomCMR><DIFTcomCMRPrimarySection>Commodities</DIFTcomCMRPrimarySection><DIFTcomCMRPrimarySectionId>MTA1-U2VjdGlvbnM=</DIFTcomCMRPrimarySectionId><DIFTcomCMRPrimaryTheme>Oil</DIFTcomCMRPrimaryTheme><DIFTcomCMRPrimaryThemeId>ZmFmYTUxOTItMGZjZC00YmJkLWJlZTQtMmY3ZDZiOWZkYmYw-VG9waWNz</DIFTcomCMRPrimaryThemeId><DIFTcomCMRBrand /><DIFTcomCMRBrandId /><DIFTcomCMRGenre>News</DIFTcomCMRGenre><DIFTcomCMRGenreId>Nw==-R2VucmVz</DIFTcomCMRGenreId><DIFTcomCMRMediaType>Text</DIFTcomCMRMediaType><DIFTcomCMRMediaTypeId>ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==</DIFTcomCMRMediaTypeId></DIFTcomCMR><DIFTcomECPositionInText>Default</DIFTcomECPositionInText><DIFTcomHideECLevel1>False</DIFTcomHideECLevel1><DIFTcomHideECLevel2>False</DIFTcomHideECLevel2><DIFTcomHideECLevel3>False</DIFTcomHideECLevel3><DIFTcomDiscussion>True</DIFTcomDiscussion><DIFTcomArticleImage>Article size</DIFTcomArticleImage></DIFTcom><DISyndication><DISyndBeenCopied>False</DISyndBeenCopied><DISyndEdition>USA</DISyndEdition><DISyndStar>01</DISyndStar><DISyndChannel /><DISyndArea /><DISyndCategory /></DISyndication></OutputChannels><EditorialNotes><Language>English</Language><Author>ravala</Author><Guides /><Editor /><Sources><Source title=\"Financial Times\"><SourceCode>FT</SourceCode><SourceDescriptor>Financial Times</SourceDescriptor><SourceOnlineInclusion>True</SourceOnlineInclusion><SourceCanBeSyndicated>True</SourceCanBeSyndicated></Source></Sources><WordCount>670</WordCount><CreationDate>20151113095342</CreationDate><EmbargoDate /><ExpiryDate>20151113095342</ExpiryDate><ObjectLocation>/FT/Content/Markets/Stories/Live/IEA MA 13.xml</ObjectLocation><OriginatingStory /><CCMS><CCMSCommissionRefNo /><CCMSContributorRefNo>CB-0000844</CCMSContributorRefNo><CCMSContributorFullName>Anjli Raval</CCMSContributorFullName><CCMSContributorInclude /><CCMSContributorRights>3</CCMSContributorRights><CCMSFilingDate /><CCMSProposedPublishingDate /></CCMS></EditorialNotes><WiresIndexing><category /><Keyword /><char_count /><priority /><basket /><title /><Version /><story_num /><file_name /><serviceid /><entry_date /><ref_field /><take_num /></WiresIndexing><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Fri Nov 13 13:29:53 GMT 2015</StartTime><EndTime>Fri Nov 13 13:29:56 GMT 2015</EndTime></ADRIS_MetaData><DFMajorCompanies><DFMajorCompany><DFCoMajScore>100</DFCoMajScore><DFCoMajDescriptor>International_Energy_Agency</DFCoMajDescriptor><DFCoMajFTCode>IEIEA00000</DFCoMajFTCode><Version>1</Version><DFCoMajTickerSymbol /><DFCoMajTickerExchangeCountry /><DFCoMajTickerExchangeCode /><DFCoMajFTMWTickercode /><DFCoMajSEDOL /><DFCoMajISIN /><DFCoMajCOFlag>O</DFCoMajCOFlag></DFMajorCompany></DFMajorCompanies><DFMinorCompanies /><DFNAICS><DFNAIC><DFNAICSFTCode>N52</DFNAICSFTCode><DFNAICSDescriptor>Finance_&amp;_Insurance</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N523</DFNAICSFTCode><DFNAICSDescriptor>Security_Commodity_Contracts_&amp;_Like_Activity</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N52321</DFNAICSFTCode><DFNAICSDescriptor>Securities_&amp;_Commodity_Exchanges</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N9261</DFNAICSFTCode><DFNAICSDescriptor>Admin_of_Economic_Programs</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N92613</DFNAICSFTCode><DFNAICSDescriptor>Regulation_&amp;_Admin_of_Utilities</DFNAICSDescriptor><Version>1997</Version></DFNAIC></DFNAICS><DFWPMIndustries /><DFFTSEGlobalClassifications /><DFStockExchangeIndices /><DFSubjects /><DFCountries /><DFRegions /><DFWPMRegions /><DFProvinces /><DFFTcomDisplayCodes /><DFFTSections /><DFWebRegions /></DataFactoryIndexing></ObjectMetadata>"
var supportedSourceCodeAttributesContentPlaceholder = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE ObjectMetadata SYSTEM \"/SysConfig/Classify/FTStories/classify.dtd\"><ObjectMetadata><EditorialDisplayIndexing><DILeadCompanies /><DITemporaryCompanies><DITemporaryCompany><DICoTempCode /><DICoTempDescriptor /><DICoTickerCode /></DITemporaryCompany></DITemporaryCompanies><DIFTSEGlobalClassifications /><DIStockExchangeIndices /><DIHotTopics /><DIHeadlineCopy>Global oil inventory stands at record level</DIHeadlineCopy><DIBylineCopy>Anjli Raval, Oil and Gas Correspondent</DIBylineCopy><DIFTNPSections /><DIFirstParCopy>International Energy Agency says crude stockpiles near 3bn barrels despite robust demand growth</DIFirstParCopy><DIMasterImgFileRef>/FT/Graphics/Online/Master_2048x1152/2015/08/New%20story%20of%20MAS_crude-oil_02.jpg?uuid=c61fcc32-61cd-11e5-9846-de406ccb37f2</DIMasterImgFileRef></EditorialDisplayIndexing><OutputChannels><DIFTN><DIFTNPublicationDate /><DIFTNZoneEdition /><DIFTNPage /><DIFTNTimeEdition /><DIFTNFronts /></DIFTN><DIFTcom><DIFTcomWebType>story</DIFTcomWebType><DIFTcomDisplayCodes><DIFTcomDisplayCodeRank1><DIFTcomDisplayCode title=\"Markets - Commodities\"><DIFTcomDisplayCodeFTCode>MKCO</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Commodities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Commodities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank1><DIFTcomDisplayCodeRank2><DIFTcomDisplayCode title=\"Markets - European Equities\"><DIFTcomDisplayCodeFTCode>MKEU</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - European Equities</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>European Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy - Oil &amp; Gas\"><DIFTcomDisplayCodeFTCode>OG8E</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy - Oil &amp; Gas</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Oil &amp; Gas</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Energy\"><DIFTcomDisplayCodeFTCode>NDEM</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Energy</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Energy</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Companies\"><DIFTcomDisplayCodeFTCode>BNIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Companies</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets\"><DIFTcomDisplayCodeFTCode>MKIP</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag /></DIFTcomDisplayCode><DIFTcomDisplayCode title=\"Markets - Equities Main\"><DIFTcomDisplayCodeFTCode>MKEQ</DIFTcomDisplayCodeFTCode><DIFTcomDisplayCodeDescriptor>Markets - Equities Main</DIFTcomDisplayCodeDescriptor><DIFTcomDisplayCodeNewsInDepth>True</DIFTcomDisplayCodeNewsInDepth><DIFTcomDisplayCodeSite>FTcom</DIFTcomDisplayCodeSite><DIFTcomDisplayCodeArticleType>News</DIFTcomDisplayCodeArticleType><DIFTcomDisplayCodeArticleBrand /><DIFTcomDisplayCodeEditorsTag>Equities</DIFTcomDisplayCodeEditorsTag></DIFTcomDisplayCode></DIFTcomDisplayCodeRank2></DIFTcomDisplayCodes><DIFTcomSubscriptionLevel>0</DIFTcomSubscriptionLevel><DIFTcomUpdateTimeStamp>False</DIFTcomUpdateTimeStamp><DIFTcomIndexAndSynd>false</DIFTcomIndexAndSynd><DIFTcomSafeToSyndicate>True</DIFTcomSafeToSyndicate><DIFTcomInitialPublication>20151113105031</DIFTcomInitialPublication><DIFTcomLastPublication>20151113132953</DIFTcomLastPublication><DIFTcomSuppresInlineAds>False</DIFTcomSuppresInlineAds><DIFTcomMap>True</DIFTcomMap><DIFTcomDisplayStyle>Normal</DIFTcomDisplayStyle><DIFTcomFeatureType>Normal</DIFTcomFeatureType><DIFTcomMarkDeleted>False</DIFTcomMarkDeleted><DIFTcomMakeUnlinkable>False</DIFTcomMakeUnlinkable><isBestStory>0</isBestStory><DIFTcomCMRId>3040370</DIFTcomCMRId><DIFTcomCMRHint /><DIFTcomCMR><DIFTcomCMRPrimarySection>Commodities</DIFTcomCMRPrimarySection><DIFTcomCMRPrimarySectionId>MTA1-U2VjdGlvbnM=</DIFTcomCMRPrimarySectionId><DIFTcomCMRPrimaryTheme>Oil</DIFTcomCMRPrimaryTheme><DIFTcomCMRPrimaryThemeId>ZmFmYTUxOTItMGZjZC00YmJkLWJlZTQtMmY3ZDZiOWZkYmYw-VG9waWNz</DIFTcomCMRPrimaryThemeId><DIFTcomCMRBrand /><DIFTcomCMRBrandId /><DIFTcomCMRGenre>News</DIFTcomCMRGenre><DIFTcomCMRGenreId>Nw==-R2VucmVz</DIFTcomCMRGenreId><DIFTcomCMRMediaType>Text</DIFTcomCMRMediaType><DIFTcomCMRMediaTypeId>ZjMwY2E2NjctMDA1Ni00ZTk4LWI0MWUtZjk5MTk2ZTMyNGVm-TWVkaWFUeXBlcw==</DIFTcomCMRMediaTypeId></DIFTcomCMR><DIFTcomECPositionInText>Default</DIFTcomECPositionInText><DIFTcomHideECLevel1>False</DIFTcomHideECLevel1><DIFTcomHideECLevel2>False</DIFTcomHideECLevel2><DIFTcomHideECLevel3>False</DIFTcomHideECLevel3><DIFTcomDiscussion>True</DIFTcomDiscussion><DIFTcomArticleImage>Article size</DIFTcomArticleImage></DIFTcom><DISyndication><DISyndBeenCopied>False</DISyndBeenCopied><DISyndEdition>USA</DISyndEdition><DISyndStar>01</DISyndStar><DISyndChannel /><DISyndArea /><DISyndCategory /></DISyndication></OutputChannels><EditorialNotes><Language>English</Language><Author>ravala</Author><Guides /><Editor /><Sources><Source title=\"Financial Times\"><SourceCode>ContentPlaceholder</SourceCode><SourceDescriptor>Financial Times</SourceDescriptor><SourceOnlineInclusion>True</SourceOnlineInclusion><SourceCanBeSyndicated>True</SourceCanBeSyndicated></Source></Sources><WordCount>670</WordCount><CreationDate>20151113095342</CreationDate><EmbargoDate /><ExpiryDate>20151113095342</ExpiryDate><ObjectLocation>/FT/Content/Markets/Stories/Live/IEA MA 13.xml</ObjectLocation><OriginatingStory /><CCMS><CCMSCommissionRefNo /><CCMSContributorRefNo>CB-0000844</CCMSContributorRefNo><CCMSContributorFullName>Anjli Raval</CCMSContributorFullName><CCMSContributorInclude /><CCMSContributorRights>3</CCMSContributorRights><CCMSFilingDate /><CCMSProposedPublishingDate /></CCMS></EditorialNotes><WiresIndexing><category /><Keyword /><char_count /><priority /><basket /><title /><Version /><story_num /><file_name /><serviceid /><entry_date /><ref_field /><take_num /></WiresIndexing><DataFactoryIndexing><ADRIS_MetaData><IndexSuccess>yes</IndexSuccess><StartTime>Fri Nov 13 13:29:53 GMT 2015</StartTime><EndTime>Fri Nov 13 13:29:56 GMT 2015</EndTime></ADRIS_MetaData><DFMajorCompanies><DFMajorCompany><DFCoMajScore>100</DFCoMajScore><DFCoMajDescriptor>International_Energy_Agency</DFCoMajDescriptor><DFCoMajFTCode>IEIEA00000</DFCoMajFTCode><Version>1</Version><DFCoMajTickerSymbol /><DFCoMajTickerExchangeCountry /><DFCoMajTickerExchangeCode /><DFCoMajFTMWTickercode /><DFCoMajSEDOL /><DFCoMajISIN /><DFCoMajCOFlag>O</DFCoMajCOFlag></DFMajorCompany></DFMajorCompanies><DFMinorCompanies /><DFNAICS><DFNAIC><DFNAICSFTCode>N52</DFNAICSFTCode><DFNAICSDescriptor>Finance_&amp;_Insurance</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N523</DFNAICSFTCode><DFNAICSDescriptor>Security_Commodity_Contracts_&amp;_Like_Activity</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N52321</DFNAICSFTCode><DFNAICSDescriptor>Securities_&amp;_Commodity_Exchanges</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N9261</DFNAICSFTCode><DFNAICSDescriptor>Admin_of_Economic_Programs</DFNAICSDescriptor><Version>1997</Version></DFNAIC><DFNAIC><DFNAICSFTCode>N92613</DFNAICSFTCode><DFNAICSDescriptor>Regulation_&amp;_Admin_of_Utilities</DFNAICSDescriptor><Version>1997</Version></DFNAIC></DFNAICS><DFWPMIndustries /><DFFTSEGlobalClassifications /><DFStockExchangeIndices /><DFSubjects /><DFCountries /><DFRegions /><DFWPMRegions /><DFProvinces /><DFFTcomDisplayCodes /><DFFTSections /><DFWebRegions /></DataFactoryIndexing></ObjectMetadata>"
