// SPDX-License-Identifier: SEE LICENSE IN LICENSE

pragma solidity ^0.8.0;

contract CrowdFunding {
    struct Campaign {
        address owner;
        string campaignType;
        string title;
        string description;
        uint256 goal;
        uint256 deadline;
        uint256 totalFunds;
        uint256 totalContributors;
        string image;
        uint256 id;
        address[] donators;
        uint256[] donations;
        bool goalReached;
        bool isDeleted;
    }

    struct Category {
        string name;
        string description;
        string image;
        uint256 id;
    }

    mapping(uint256 => Campaign) public campaigns;
    mapping(uint256 => Category) public categories;

    uint public campaignCount = 0;
    uint public categoryCount = 0;

    modifier isTimePassed(uint256 _campaignId) {
        require(
            block.timestamp > campaigns[_campaignId].deadline,
            "Deadline is not passed"
        );
        _;
    }

    modifier isCampaignOwner(uint256 _campaignId) {
        require(
            msg.sender == campaigns[_campaignId].owner,
            "You are not the owner of this campaign"
        );
        _;
    }

    //? Create a new campaign
    function createCampaign(
        string memory _campaignType,
        string memory _title,
        string memory _description,
        uint256 _goal,
        uint256 _deadline,
        string memory _image
    ) public returns (uint256) {
        Campaign storage campaign = campaigns[campaignCount];

        require(campaign.deadline < block.timestamp, "Deadline is not passed");

        campaign.owner = msg.sender;
        campaign.campaignType = _campaignType;
        campaign.title = _title;
        campaign.description = _description;
        campaign.goal = _goal;
        campaign.deadline = _deadline;
        campaign.image = _image;
        campaign.id = campaignCount;
        campaign.isDeleted = false;
        campaign.goalReached = false;
        campaignCount++;

        return campaignCount - 1;
    }

    //? Get campaign by id
    function getCampaign(
        uint256 _campaignId
    ) public view returns (Campaign memory) {
        return campaigns[_campaignId];
    }

    //? Donate to a campaign
    function donate(uint256 _campaignId) public payable {
        uint256 amount = msg.value;
        Campaign storage campaign = campaigns[_campaignId];

        require(campaign.deadline > block.timestamp, "Deadline is passed");
        require(msg.value > 0, "Donation amount must be greater than 0");
        require(campaign.isDeleted == false, "Campaign is deleted");
        if (campaign.goalReached) {
            revert("Campaign goal is reached, you can't donate");
        }

        campaign.totalFunds += msg.value;
        campaign.totalContributors += 1;
        campaign.donators.push(msg.sender);
        campaign.donations.push(msg.value);
        campaign.id = _campaignId;
        campaign.owner = campaign.owner;

        // check if campaign goal is reached
        if (campaign.totalFunds + msg.value >= campaign.goal) {
            campaign.goalReached = true;
        }

        // make sure the amount plus the totalFunds is less than the goal amount if not refund the difference
        if (campaign.totalFunds > campaign.goal) {
            uint256 difference = campaign.totalFunds - campaign.goal;
            (bool sent, ) = payable(msg.sender).call{value: difference}("");
            if (sent) {
                campaign.totalFunds -= difference;
            } else {
                revert("Failed to send Ether");
            }
        }

        // pay donation to campaign owner
        if (campaign.goalReached) {
            (bool sent, ) = payable(campaign.owner).call{value: amount}("");
            if (sent) {
                campaign.totalFunds += amount;
            } else {
                revert("Failed to send Ether");
            }
        }

        // if goals is reached set goalReached to true
        if (campaign.totalFunds >= campaign.goal) {
            campaign.goalReached = true;
        }
    }

    //? Get campaign donations
    function getCampaignDonators(
        uint256 _campaignId
    ) public view returns (address[] memory) {
        return campaigns[_campaignId].donators;
    }

    //? Get campaign donations
    function getCampaignDonations(
        uint256 _campaignId
    ) public view returns (uint256[] memory) {
        return campaigns[_campaignId].donations;
    }

    //? Get campaign donations sum
    function getCampaignDonationsSum(
        uint256 _campaignId
    ) public view returns (uint256) {
        uint256 sum = 0;

        for (uint256 i = 0; i < campaigns[_campaignId].donations.length; i++) {
            sum += campaigns[_campaignId].donations[i];
        }

        return sum;
    }

    //? Get campaign donations count
    function getCampaignDonationsCount(
        uint256 _campaignId
    ) public view returns (uint256) {
        return campaigns[_campaignId].donations.length;
    }

    //? Get campaign donators count
    function getCampaignDonatorsCount(
        uint256 _campaignId
    ) public view returns (uint256) {
        return campaigns[_campaignId].donators.length;
    }

    //? Get campaign donations count
    function getCampaignsCount() public view returns (uint256) {
        return campaignCount;
    }

    //? Get total donations by campaign id
    function getTotalDonationsByCampaignId(
        uint256 _campaignId
    ) public view returns (uint256) {
        // number of donations
        uint256 donationsCount = campaigns[_campaignId].donations.length;
        return donationsCount;
    }

    //? Get campaign donations count
    function getCampaigns() public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        for (uint256 i = 0; i < campaignCount; i++) {
            _campaigns[i] = campaigns[i];
        }

        return _campaigns;
    }

    //? Get campaign donations count
    function getCampaignsByOwner(
        address _owner
    ) public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (campaigns[i].owner == _owner) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return _campaigns;
    }

    //? Get campaign By Type
    function getCampaignsByType(
        uint256 _categoryId
    ) public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (
                keccak256(abi.encodePacked(campaigns[i].campaignType)) ==
                keccak256(abi.encodePacked(categories[_categoryId].name))
            ) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return _campaigns;
    }

    //? Get Donors Addresses And Amounts
    function getDonorsAddressesAndAmounts(
        uint256 _campaignId
    ) public view returns (address[] memory, uint256[] memory, uint256) {
        return (
            campaigns[_campaignId].donators,
            campaigns[_campaignId].donations,
            campaigns[_campaignId].totalFunds
        );
    }

    //? Pay out to campaign owner if goal is reached
    function payOut(
        uint256 _campaignId
    ) public isCampaignOwner(_campaignId) isTimePassed(_campaignId) {
        Campaign storage campaign = campaigns[_campaignId];

        require(
            campaign.totalFunds >= campaign.goal,
            "Campaign goal is not reached"
        );

        (bool sent, ) = payable(campaign.owner).call{
            value: campaign.totalFunds
        }("");

        if (sent) {
            campaign.totalFunds = 0;
        } else {
            revert("Failed to send Ether");
        }
    }

    //? Send back donations if goal is not reached
    function sendBackDonations(
        uint256 _campaignId
    ) public isCampaignOwner(_campaignId) isTimePassed(_campaignId) {
        Campaign storage campaign = campaigns[_campaignId];

        require(
            campaign.totalFunds < campaign.goal,
            "Campaign goal is reached"
        );

        for (uint256 i = 0; i < campaign.donators.length; i++) {
            (bool sent, ) = payable(campaign.donators[i]).call{
                value: campaign.donations[i]
            }("");

            if (sent) {
                campaign.totalFunds -= campaign.donations[i];
            } else {
                revert("Failed to send Ether");
            }
        }
    }

    function createCategory(
        string memory _name,
        string memory _description,
        string memory _image
    ) public returns (uint256) {
        Category storage category = categories[categoryCount];

        category.name = _name;
        category.description = _description;
        category.image = _image;
        category.id = categoryCount;

        categoryCount++;

        return categoryCount - 1;
    }

    //? Get all catrgories
    function getCategories() public view returns (Category[] memory) {
        Category[] memory _categories = new Category[](categoryCount);

        for (uint256 i = 0; i < categoryCount; i++) {
            _categories[i] = categories[i];
        }

        return _categories;
    }

    //? Search Campaign by title
    function searchCampaignByName(
        string memory _name
    ) public view returns (Campaign[] memory) {
        string[] memory keywords = splitString(_name, " ");

        Campaign[] memory _campaigns = new Campaign[](campaignCount);
        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (matchesKeywords(campaigns[i], keywords)) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return resizeCampaignArray(_campaigns, count);
    }

    function resizeCampaignArray(
        Campaign[] memory _array,
        uint256 _size
    ) internal pure returns (Campaign[] memory) {
        Campaign[] memory result = new Campaign[](_size);
        for (uint256 i = 0; i < _size; i++) {
            result[i] = _array[i];
        }
        return result;
    }

    function splitString(
        string memory _text,
        string memory _delimiter
    ) internal pure returns (string[] memory) {
        bytes memory bytesText = bytes(_text);
        bytes memory bytesDelimiter = bytes(_delimiter);

        uint256 count = 1;
        for (uint256 i = 0; i < bytesText.length; i++) {
            if (bytesText[i] == bytesDelimiter[0]) {
                count++;
            }
        }

        string[] memory parts = new string[](count);
        uint256 startIndex = 0;
        count = 0;
        for (uint256 i = 0; i < bytesText.length; i++) {
            if (bytesText[i] == bytesDelimiter[0]) {
                parts[count] = string(bytesSubstring(bytesText, startIndex, i));
                startIndex = i + 1;
                count++;
            }
        }
        parts[count] = string(
            bytesSubstring(bytesText, startIndex, bytesText.length)
        );

        return parts;
    }

    function bytesSubstring(
        bytes memory _bytes,
        uint256 _start,
        uint256 _length
    ) internal pure returns (bytes memory) {
        require(_start + _length <= _bytes.length, "Invalid substring range");
        bytes memory result = new bytes(_length);
        for (uint256 i = 0; i < _length; i++) {
            result[i] = _bytes[_start + i];
        }
        return result;
    }

    function matchesKeywords(
        Campaign memory _campaign,
        string[] memory _keywords
    ) internal pure returns (bool) {
        for (uint256 i = 0; i < _keywords.length; i++) {
            if (
                containsKeyword(_campaign.title, _keywords[i]) ||
                containsKeyword(_campaign.description, _keywords[i]) ||
                containsKeyword(_campaign.campaignType, _keywords[i])
            ) {
                return true;
            }
        }
        return false;
    }

    function containsKeyword(
        string memory _text,
        string memory _keyword
    ) internal pure returns (bool) {
        return
            (bytes(_text).length > 0) &&
            (bytes(_keyword).length > 0) &&
            (bytes(_text).length >= bytes(_keyword).length) &&
            (keccak256(abi.encodePacked(_text)) ==
                keccak256(abi.encodePacked(_keyword)));
    }

    function deleteCampaign(
        uint256 _campaignId
    ) public isCampaignOwner(_campaignId) {
        campaigns[_campaignId].isDeleted = true;
    }

    function deleteCategory(uint256 _categoryId) public {
        delete categories[_categoryId];
    }

    function deleteAllCampaigns() public {
        for (uint256 i = 0; i < campaignCount; i++) {
            // set isDeleted to true
            campaigns[i].isDeleted = true;
        }
    }

    function deleteAllCategories() public {
        for (uint256 i = 0; i < categoryCount; i++) {
            delete categories[i];
        }
    }

    // get Campaigns by category
    function getCampaignsByCategory(
        uint256 _categoryId
    ) public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (
                keccak256(abi.encodePacked(campaigns[i].campaignType)) ==
                keccak256(abi.encodePacked(categories[_categoryId].name))
            ) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return _campaigns;
    }
}
